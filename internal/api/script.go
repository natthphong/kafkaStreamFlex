package api

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/natthphong/kafkaStreamFlex/config"
	"github.com/natthphong/kafkaStreamFlex/internal/connection/logz"
	"github.com/natthphong/kafkaStreamFlex/internal/connection/s3"
	"github.com/natthphong/kafkaStreamFlex/internal/models"
	"github.com/natthphong/kafkaStreamFlex/internal/repository"
	"github.com/natthphong/kafkaStreamFlex/internal/script"
	"github.com/natthphong/kafkaStreamFlex/pkg/utils"
	"go.uber.org/zap"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func NewUploadScript(
	config config.ScriptConfig,
	scriptRepository repository.ScriptRepository,
	s3UploadFunc s3.S3UploadFunc,
) fiber.Handler {
	return func(c *fiber.Ctx) error {
		logger := logz.NewLogger()
		ctx := c.Context()
		id := c.FormValue("id", "")
		if id == "" {
			return BadRequest(c, "id is required")
		}
		fileHeader, err := c.FormFile("file")
		if err != nil {
			return BadRequest(c, "file is required")
		}
		basePath := config.BasePath
		if strings.HasSuffix(basePath, "/") {
			basePath = strings.TrimSuffix(basePath, "/")
		}
		originalFileName := strings.Split(fileHeader.Filename, ".")
		fileName := fileHeader.Filename
		expression := ""
		if len(originalFileName) > 1 {
			fileName = originalFileName[0]
			expression = originalFileName[1]
		}
		if expression != "zip" {
			return BadRequest(c, "invalid expression")
		}

		ref := uuid.NewString()
		current := time.Now()
		currentMonth := current.Format("2200601")
		tempUpload := "./upload"
		uploadZipDir := fmt.Sprintf("%s/zip/%s/%s", tempUpload, fileName, currentMonth)
		uploadScriptDir := fmt.Sprintf("%s/script/%s/%s/%s", tempUpload, fileName, currentMonth, ref)
		uploadSoFileDir := fmt.Sprintf("%s/so/%s/%s", basePath, fileName, currentMonth)
		err = os.MkdirAll(uploadZipDir, os.ModePerm)
		if err != nil {
			return err
		}
		err = os.MkdirAll(uploadScriptDir, os.ModePerm)
		if err != nil {
			return err
		}
		err = os.MkdirAll(uploadSoFileDir, os.ModePerm)
		if err != nil {
			return err
		}
		saveZipFilePath := filepath.Join(uploadZipDir, ref+".zip")
		err = c.SaveFile(fileHeader, saveZipFilePath)
		if err != nil {
			logger.Error(err.Error())
			return fiber.NewError(fiber.StatusInternalServerError, "Failed to save file")
		}

		err = utils.UnzipFile(saveZipFilePath, uploadScriptDir)
		if err != nil {
			logger.Error(err.Error())
			return fiber.NewError(fiber.StatusInternalServerError, "Failed to unzip file: "+err.Error())
		}
		saveSoFilePath := filepath.Join(uploadSoFileDir, ref+".so")
		err = script.BuildPlugin(uploadScriptDir, saveSoFilePath)
		if err != nil {
			logger.Info(uploadScriptDir, zap.String("saveSoFilePath", saveSoFilePath))
			logger.Error(err.Error())
			return InternalError(c, "Failed to build plugin")
		}

		soBytes, err := os.ReadFile(saveSoFilePath)
		if err != nil {
			return InternalError(c, "Failed to read .so file")
		}

		err = s3UploadFunc(logger, saveSoFilePath, soBytes)
		if err != nil {
			return InternalError(c, err.Error())
		}

		scriptEntity := models.Script{
			ID:         id,
			ScriptKey:  saveSoFilePath,
			ScriptName: ref,
		}
		err = scriptRepository.Save(ctx, logger, scriptEntity)
		if err != nil {
			return err
		}
		return Ok(c, fiber.Map{
			"path": saveSoFilePath,
			"ref":  ref,
		})
	}
}

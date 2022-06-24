package exercise

import (
	"course/internal/domain"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ExerciseService struct {
	db *gorm.DB
}

func NewExerciseService(database *gorm.DB) *ExerciseService {
	return &ExerciseService{
		db: database,
	}
}

func (ex ExerciseService) GetExercise(ctx *gin.Context) {
	paramID := ctx.Param("id")
	id, err := strconv.Atoi(paramID)
	if err != nil {
		ctx.JSON(400, gin.H{
			"message": "invalid exercise id",
		})
		return
	}
	var exercise domain.Exercise
	err = ex.db.Where("id = ?", id).Preload("Questions").Take(&exercise).Error
	if err != nil {
		ctx.JSON(404, gin.H{
			"message": "not found",
		})
		return
	}
	ctx.JSON(200, exercise)
}

func (ex ExerciseService) GetUserScore(ctx *gin.Context) {
	paramExerciseID := ctx.Param("id")
	exerciseID, err := strconv.Atoi(paramExerciseID)
	if err != nil {
		ctx.JSON(400, gin.H{
			"message": "invalid exercise id",
		})
		return
	}
	var exercise domain.Exercise
	err = ex.db.Where("id = ?", exerciseID).Preload("Questions").Take(&exercise).Error
	if err != nil {
		ctx.JSON(404, gin.H{
			"message": "not found",
		})
		return
	}

	userID := int(ctx.Request.Context().Value("user_id").(float64))
	var answers []domain.Answer
	err = ex.db.Where("exercise_id = ? AND user_id = ?", exerciseID, userID).Find(&answers).Error

	if err != nil {
		ctx.JSON(200, gin.H{
			"score": 0,
		})
		return
	}

	mapQA := make(map[int]domain.Answer)
	for _, answer := range answers {
		mapQA[answer.QuestionID] = answer
	}

	var score int
	for _, question := range exercise.Questions {
		if strings.EqualFold(question.CorrectAnswer, mapQA[question.ID].Answer) {
			score += question.Score
		}
	}
	ctx.JSON(200, gin.H{
		"score": score,
	})
}

func (ex ExerciseService) CreateExcercise(ctx *gin.Context) {
	var excercise domain.Exercise
	err := ctx.ShouldBind(&excercise)

	if err != nil {
		ctx.JSON(400, gin.H{
			"message": "invalid input",
		})

		return
	}

	if excercise.Title == "" {
		ctx.JSON(400, gin.H{
			"message": "field title required",
		})

		return
	}

	if excercise.Description == "" {
		ctx.JSON(400, gin.H{
			"message": "field description required",
		})

		return
	}

	if err :=
		ex.db.Create(&excercise).Error; err != nil {
		ctx.JSON(500, gin.H{
			"message": "failed when creating excercise",
		})
		return
	}

	ctx.JSON(201, gin.H{
		"result": &excercise,
	})
}

func (ex ExerciseService) CreateQuestion(ctx *gin.Context) {

	var question domain.Question

	question.CreatorID = int(ctx.Request.Context().Value("user_id").(float64))

	err := ctx.ShouldBind(&question)

	if err != nil {
		ctx.JSON(400, gin.H{
			"message": "invalid input",
		})

		return
	}

	if question.ExerciseID == 0 {
		ctx.JSON(400, gin.H{
			"message": "field ExerciseID is required",
		})

		return
	}
	if question.Body == "" {
		ctx.JSON(400, gin.H{
			"message": "field Body is required",
		})

		return
	}
	if question.OptionA == "" {
		ctx.JSON(400, gin.H{
			"message": "field OptionA is required",
		})

		return
	}
	if question.OptionB == "" {
		ctx.JSON(400, gin.H{
			"message": "field OptionB is required",
		})

		return
	}
	if question.OptionC == "" {
		ctx.JSON(400, gin.H{
			"message": "field OptionC is required",
		})

		return
	}
	if question.OptionD == "" {
		ctx.JSON(400, gin.H{
			"message": "field OptionD is required",
		})

		return
	}
	if question.CorrectAnswer == "" {
		ctx.JSON(400, gin.H{
			"message": "field CorrectAnswer is required",
		})

		return
	}
	if question.Score == 0 {
		ctx.JSON(400, gin.H{
			"message": "field Score is required",
		})

		return
	}

	if err := ex.db.Create(&question).Error; err != nil {
		ctx.JSON(500, gin.H{
			"message": "failed when creating question",
		})

		return
	}

	ctx.JSON(201, gin.H{
		"result": &question,
	})
}

func (ex *ExerciseService) CreateAnswer(ctx *gin.Context) {
	var answer domain.Answer

	answer.UserID = int(ctx.Request.Context().Value("user_id").(float64))

	err := ctx.ShouldBind(&answer)

	if err != nil {
		ctx.JSON(400, gin.H{
			"message": "invalid input",
		})

		return
	}

	if answer.Answer == "" {
		ctx.JSON(400, gin.H{
			"message": "field answer required",
		})

		return
	}

	if answer.QuestionID == 0 {
		ctx.JSON(400, gin.H{
			"message": "field question_id required",
		})

		return
	}

	if answer.ExerciseID == 0 {
		ctx.JSON(400, gin.H{
			"message": "field excecise_id required",
		})

		return
	}

	if err :=
		ex.db.Create(&answer).Error; err != nil {
		ctx.JSON(500, gin.H{
			"message": "failed when creating answer ",
		})

		return
	}

	ctx.JSON(201, gin.H{
		"result": &answer,
	})
}

package recipe

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"github.com/yuin/goldmark"
	"gluttony/internal/ingredient"
	"gluttony/internal/recipe/queries"
	"io"
	"time"
)

type Tag struct {
	ID    int
	Order int
	Name  string
}

type CreateInput struct {
	Name            string
	Description     string
	Source          string
	Instructions    string
	Servings        int8
	PreparationTime time.Duration
	CookTime        time.Duration
	Tags            []string
	Ingredients     []Ingredient
	Nutrition       Nutrition
	ThumbnailImage  io.Reader
}

type Service struct {
	db         *sql.DB
	queries    *queries.Queries
	mediaStore MediaStore
	markdown   goldmark.Markdown
}

func NewService(db *sql.DB, mediaStore MediaStore) *Service {
	if db == nil {
		panic("db is nil")
	}

	if mediaStore == nil {
		panic("mediaStore is nil")
	}

	return &Service{
		queries:    queries.New(db),
		db:         db,
		mediaStore: mediaStore,
		markdown:   goldmark.New(),
	}
}

func (s *Service) Create(ctx context.Context, input CreateInput) (err error) {
	thumbnailImageURL := ""
	if input.ThumbnailImage != nil {
		thumbnailImageURL, err = s.mediaStore.UploadImage(input.ThumbnailImage)
		if err != nil {
			return fmt.Errorf("upload thumbnail image: %w", err)
		}
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err == nil {
			err = tx.Commit()
			return
		}

		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			err = fmt.Errorf("recipe create: %w: %w", err, rollbackErr)

			// TODO: remove image
		}
	}()

	txQueries := s.queries.WithTx(tx)

	ingredients, err := s.createOrGetIngredients(ctx, txQueries, input.Ingredients)
	if err != nil {
		return fmt.Errorf("create ingredients: %w", err)
	}

	tags, err := s.createOrGetTags(ctx, txQueries, input.Tags)
	if err != nil {
		return fmt.Errorf("create tags: %w", err)
	}

	createRecipeParams := queries.CreateRecipeParams{
		Name:                   input.Name,
		Description:            input.Description,
		InstructionsMarkdown:   input.Instructions,
		CookTimeSeconds:        int64(input.CookTime.Seconds()),
		PreparationTimeSeconds: int64(input.PreparationTime.Seconds()),
		Source:                 input.Source,
	}
	if thumbnailImageURL != "" {
		createRecipeParams.ThumbnailUrl = sql.NullString{
			String: thumbnailImageURL,
			Valid:  true,
		}
	}

	recipeID, err := txQueries.CreateRecipe(ctx, createRecipeParams)
	if err != nil {
		return fmt.Errorf("create recipe: %w", err)
	}

	err = txQueries.CreateNutrition(ctx, queries.CreateNutritionParams{
		RecipeID: recipeID,
		Calories: float64(input.Nutrition.Calories),
		Fat:      float64(input.Nutrition.Fat),
		Carbs:    float64(input.Nutrition.Carbs),
		Protein:  float64(input.Nutrition.Protein),
	})
	if err != nil {
		return fmt.Errorf("create nutrition: %w", err)
	}

	for i := range tags {
		err = txQueries.CreateRecipeTag(ctx, queries.CreateRecipeTagParams{
			RecipeOrder: int64(tags[i].Order),
			RecipeID:    recipeID,
			TagID:       int64(tags[i].ID),
		})
		if err != nil {
			return fmt.Errorf("create recipe tag: %w", err)
		}
	}

	for i := range ingredients {
		err = txQueries.CreateRecipeIngredient(ctx, queries.CreateRecipeIngredientParams{
			RecipeOrder:  int64(ingredients[i].Order),
			RecipeID:     recipeID,
			IngredientID: int64(ingredients[i].ID),
			Unit:         ingredients[i].Unit,
			Quantity:     int64(ingredients[i].Quantity),
		})
	}

	return nil
}

func (s *Service) createOrGetTags(
	ctx context.Context,
	txQueries *queries.Queries,
	tags []string,
) ([]Tag, error) {

	existing, err := txQueries.AllTagsByNames(ctx, tags)
	if err != nil {
		return nil, fmt.Errorf("get ingredients: %w", err)
	}

	existingLookup := make(map[string]queries.Tag, len(existing))
	for i := range existing {
		existingLookup[existing[i].Name] = existing[i]
	}

	out := make([]Tag, len(tags))
	for i := range tags {
		var tagID int64
		if value, ok := existingLookup[tags[i]]; ok {
			tagID = value.ID
		} else {
			id, err := txQueries.CreateTag(ctx, tags[i])
			if err != nil {
				return nil, fmt.Errorf("create tag: %w", err)
			}

			tagID = id
		}

		out[i].Name = tags[i]
		out[i].ID = int(tagID)
		out[i].Order = i
	}

	return out, nil
}

func (s *Service) createOrGetIngredients(
	ctx context.Context,
	txQueries *queries.Queries,
	ingredients []Ingredient,
) ([]Ingredient, error) {

	names := make([]string, len(ingredients))
	for i := range ingredients {
		names[i] = ingredients[i].Name
	}

	existingIngredients, err := txQueries.AllIngredientsByNames(ctx, names)
	if err != nil {
		return nil, fmt.Errorf("get ingredients: %w", err)
	}

	existingLookup := make(map[string]queries.Ingredient, len(existingIngredients))
	for i := range existingIngredients {
		existingLookup[existingIngredients[i].Name] = existingIngredients[i]
	}

	out := make([]Ingredient, len(ingredients))
	for i := range ingredients {
		var ingredientID int64
		if value, ok := existingLookup[ingredients[i].Name]; ok {
			ingredientID = value.ID
		} else {
			id, err := txQueries.CreateIngredient(ctx, ingredients[i].Name)
			if err != nil {
				return nil, fmt.Errorf("create ingredient: %w", err)
			}

			ingredientID = id
		}

		out[i].ID = int(ingredientID)
		out[i].Name = ingredients[i].Name
		out[i].Order = int8(i)
		out[i].Quantity = ingredients[i].Quantity
		out[i].Unit = ingredients[i].Unit
	}

	return out, nil
}

func (s *Service) AllPartial(ctx context.Context, input SearchInput) ([]Partial, error) {
	recipes, err := s.queries.AllPartialRecipes(ctx, sql.NullString{
		String: input.Query,
		Valid:  input.Query != "",
	})
	if err != nil {
		return []Partial{}, fmt.Errorf("get all recipes: %w", err)
	}

	recipeIDs := make([]int64, len(recipes))
	out := make([]Partial, len(recipes))
	for i := range recipes {
		recipeIDs[i] = recipes[i].ID

		out[i].ID = int(recipes[i].ID)
		out[i].Name = recipes[i].Name
		out[i].Description = recipes[i].Description

		if recipes[i].ThumbnailUrl.Valid {
			out[i].ThumbnailImageURL = recipes[i].ThumbnailUrl.String
		}
	}

	tags, err := s.allTagsByRecipeIDs(ctx, recipeIDs)
	if err != nil {
		return []Partial{}, err
	}

	for i := range out {
		recipeTags, ok := tags[int64(out[i].ID)]
		if !ok {
			continue
		}

		out[i].Tags = recipeTags
	}

	return out, nil
}

func (s *Service) GetFull(ctx context.Context, recipeID int64) (Full, error) {
	recipe, err := s.queries.GetFullRecipe(ctx, recipeID)
	if err != nil {
		return Full{}, fmt.Errorf("get recipe id=%d: %w", recipeID, err)
	}

	tags, err := s.allTagsByRecipeIDs(ctx, []int64{recipeID})
	if err != nil {
		return Full{}, fmt.Errorf("get all recipe tags for recipe id=%d: %w", recipeID, err)
	}

	ingredients, err := s.allIngredientsByRecipeIDs(ctx, []int64{recipeID})
	if err != nil {
		return Full{}, fmt.Errorf("get all ingredients for recipe id=%d: %w", recipeID, err)
	}

	out := Full{
		ID:                int(recipeID),
		Name:              recipe.Name,
		Description:       recipe.Description,
		ThumbnailImageURL: recipe.ThumbnailUrl.String,
		Tags:              tags[recipeID],
		Source:            recipe.Source,
		Servings:          int8(recipe.Servings),
		PreparationTime:   time.Duration(recipe.PreparationTimeSeconds) * time.Second,
		CookTime:          time.Duration(recipe.CookTimeSeconds) * time.Second,
		Ingredients:       ingredients[recipeID],
		Nutrition: Nutrition{
			Calories: float32(recipe.Calories),
			Fat:      float32(recipe.Fat),
			Carbs:    float32(recipe.Carbs),
			Protein:  float32(recipe.Protein),
		},
	}

	var htmlInstructions bytes.Buffer
	if err := s.markdown.Convert([]byte(recipe.InstructionsMarkdown), &htmlInstructions); err != nil {
		return Full{}, fmt.Errorf("convert markdown instructions to html: %w", err)
	}

	out.Instructions = htmlInstructions.String()

	return out, nil
}

func (s *Service) allIngredientsByRecipeIDs(
	ctx context.Context,
	recipeIDs []int64,
) (map[int64][]Ingredient, error) {
	ingredients, err := s.queries.AllRecipeIngredients(ctx, recipeIDs)
	if err != nil {
		return nil, fmt.Errorf("get all ingredients by recipe ids = %+v: %w", recipeIDs, err)
	}

	out := make(map[int64][]Ingredient, len(ingredients))
	for i := range ingredients {
		if out[ingredients[i].RecipeID] == nil {
			out[ingredients[i].RecipeID] = []Ingredient{}
		}

		out[ingredients[i].RecipeID] = append(out[ingredients[i].RecipeID], Ingredient{
			Ingredient: ingredient.Ingredient{
				ID:   int(ingredients[i].ID),
				Name: ingredients[i].Name,
			},
			Order:    int8(ingredients[i].RecipeOrder),
			Quantity: float32(ingredients[i].Quantity),
			Unit:     ingredients[i].Unit,
		})
	}

	return out, nil
}

func (s *Service) allTagsByRecipeIDs(ctx context.Context, recipeIDs []int64) (map[int64][]Tag, error) {
	tags, err := s.queries.AllRecipeTags(ctx, recipeIDs)
	if err != nil {
		return nil, fmt.Errorf("get all tags by recipe ids = %+v: %w", recipeIDs, err)
	}

	out := make(map[int64][]Tag, len(tags))
	for i := range tags {
		if out[tags[i].RecipeID] == nil {
			out[tags[i].RecipeID] = []Tag{}
		}

		out[tags[i].RecipeID] = append(out[tags[i].RecipeID], Tag{
			ID:    int(tags[i].ID),
			Order: int(tags[i].RecipeOrder),
			Name:  tags[i].Name,
		})
	}

	return out, nil
}

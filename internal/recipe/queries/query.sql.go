// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: query.sql

package queries

import (
	"context"
	"database/sql"
	"strings"
	"time"
)

const allIngredientsByNames = `-- name: AllIngredientsByNames :many
SELECT id, name
FROM ingredients
WHERE name in (/*SLICE:names*/?)
`

func (q *Queries) AllIngredientsByNames(ctx context.Context, names []string) ([]Ingredient, error) {
	query := allIngredientsByNames
	var queryParams []interface{}
	if len(names) > 0 {
		for _, v := range names {
			queryParams = append(queryParams, v)
		}
		query = strings.Replace(query, "/*SLICE:names*/?", strings.Repeat(",?", len(names))[1:], 1)
	} else {
		query = strings.Replace(query, "/*SLICE:names*/?", "NULL", 1)
	}
	rows, err := q.db.QueryContext(ctx, query, queryParams...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Ingredient
	for rows.Next() {
		var i Ingredient
		if err := rows.Scan(&i.ID, &i.Name); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const allPartialRecipes = `-- name: AllPartialRecipes :many
SELECT id, name, description, thumbnail_url
FROM recipes
ORDER BY id desc
`

type AllPartialRecipesRow struct {
	ID           int64
	Name         string
	Description  string
	ThumbnailUrl sql.NullString
}

func (q *Queries) AllPartialRecipes(ctx context.Context) ([]AllPartialRecipesRow, error) {
	rows, err := q.db.QueryContext(ctx, allPartialRecipes)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []AllPartialRecipesRow
	for rows.Next() {
		var i AllPartialRecipesRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Description,
			&i.ThumbnailUrl,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const allRecipeIngredients = `-- name: AllRecipeIngredients :many
SELECT id, name, recipe_order, recipe_id, ingredient_id, unit, quantity
FROM ingredients
         JOIN recipe_ingredients ri on ingredients.id = ri.ingredient_id
WHERE recipe_id in (/*SLICE:ids*/?)
ORDER BY recipe_id, recipe_order
`

type AllRecipeIngredientsRow struct {
	ID           int64
	Name         string
	RecipeOrder  int64
	RecipeID     int64
	IngredientID int64
	Unit         string
	Quantity     int64
}

func (q *Queries) AllRecipeIngredients(ctx context.Context, ids []int64) ([]AllRecipeIngredientsRow, error) {
	query := allRecipeIngredients
	var queryParams []interface{}
	if len(ids) > 0 {
		for _, v := range ids {
			queryParams = append(queryParams, v)
		}
		query = strings.Replace(query, "/*SLICE:ids*/?", strings.Repeat(",?", len(ids))[1:], 1)
	} else {
		query = strings.Replace(query, "/*SLICE:ids*/?", "NULL", 1)
	}
	rows, err := q.db.QueryContext(ctx, query, queryParams...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []AllRecipeIngredientsRow
	for rows.Next() {
		var i AllRecipeIngredientsRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.RecipeOrder,
			&i.RecipeID,
			&i.IngredientID,
			&i.Unit,
			&i.Quantity,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const allRecipeSummaryByIDs = `-- name: AllRecipeSummaryByIDs :many
SELECT id, name, description, thumbnail_url
FROM recipes
WHERE id in (/*SLICE:ids*/?)
`

type AllRecipeSummaryByIDsRow struct {
	ID           int64
	Name         string
	Description  string
	ThumbnailUrl sql.NullString
}

func (q *Queries) AllRecipeSummaryByIDs(ctx context.Context, ids []int64) ([]AllRecipeSummaryByIDsRow, error) {
	query := allRecipeSummaryByIDs
	var queryParams []interface{}
	if len(ids) > 0 {
		for _, v := range ids {
			queryParams = append(queryParams, v)
		}
		query = strings.Replace(query, "/*SLICE:ids*/?", strings.Repeat(",?", len(ids))[1:], 1)
	} else {
		query = strings.Replace(query, "/*SLICE:ids*/?", "NULL", 1)
	}
	rows, err := q.db.QueryContext(ctx, query, queryParams...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []AllRecipeSummaryByIDsRow
	for rows.Next() {
		var i AllRecipeSummaryByIDsRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Description,
			&i.ThumbnailUrl,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const allRecipeTags = `-- name: AllRecipeTags :many
SELECT id, name, recipe_order, recipe_id, tag_id
FROM tags
         JOIN recipe_tags rt on tags.id = rt.tag_id
WHERE recipe_id in (/*SLICE:ids*/?)
ORDER BY recipe_id, recipe_order
`

type AllRecipeTagsRow struct {
	ID          int64
	Name        string
	RecipeOrder int64
	RecipeID    int64
	TagID       int64
}

func (q *Queries) AllRecipeTags(ctx context.Context, ids []int64) ([]AllRecipeTagsRow, error) {
	query := allRecipeTags
	var queryParams []interface{}
	if len(ids) > 0 {
		for _, v := range ids {
			queryParams = append(queryParams, v)
		}
		query = strings.Replace(query, "/*SLICE:ids*/?", strings.Repeat(",?", len(ids))[1:], 1)
	} else {
		query = strings.Replace(query, "/*SLICE:ids*/?", "NULL", 1)
	}
	rows, err := q.db.QueryContext(ctx, query, queryParams...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []AllRecipeTagsRow
	for rows.Next() {
		var i AllRecipeTagsRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.RecipeOrder,
			&i.RecipeID,
			&i.TagID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const allTagsByNames = `-- name: AllTagsByNames :many
SELECT id, name
FROM tags
WHERE name in (/*SLICE:names*/?)
`

func (q *Queries) AllTagsByNames(ctx context.Context, names []string) ([]Tag, error) {
	query := allTagsByNames
	var queryParams []interface{}
	if len(names) > 0 {
		for _, v := range names {
			queryParams = append(queryParams, v)
		}
		query = strings.Replace(query, "/*SLICE:names*/?", strings.Repeat(",?", len(names))[1:], 1)
	} else {
		query = strings.Replace(query, "/*SLICE:names*/?", "NULL", 1)
	}
	rows, err := q.db.QueryContext(ctx, query, queryParams...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Tag
	for rows.Next() {
		var i Tag
		if err := rows.Scan(&i.ID, &i.Name); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const createIngredient = `-- name: CreateIngredient :one
INSERT INTO ingredients (name)
VALUES (?)
RETURNING id
`

func (q *Queries) CreateIngredient(ctx context.Context, name string) (int64, error) {
	row := q.db.QueryRowContext(ctx, createIngredient, name)
	var id int64
	err := row.Scan(&id)
	return id, err
}

const createNutrition = `-- name: CreateNutrition :exec
INSERT INTO recipe_nutrition (recipe_id, calories, fat, carbs, protein)
VALUES (?, ?, ?, ?, ?)
`

type CreateNutritionParams struct {
	RecipeID int64
	Calories float64
	Fat      float64
	Carbs    float64
	Protein  float64
}

func (q *Queries) CreateNutrition(ctx context.Context, arg CreateNutritionParams) error {
	_, err := q.db.ExecContext(ctx, createNutrition,
		arg.RecipeID,
		arg.Calories,
		arg.Fat,
		arg.Carbs,
		arg.Protein,
	)
	return err
}

const createRecipe = `-- name: CreateRecipe :one
INSERT INTO recipes (name, description, instructions_markdown, thumbnail_url,
                     cook_time_seconds, preparation_time_seconds, source, updated_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?)
RETURNING id
`

type CreateRecipeParams struct {
	Name                   string
	Description            string
	InstructionsMarkdown   string
	ThumbnailUrl           sql.NullString
	CookTimeSeconds        int64
	PreparationTimeSeconds int64
	Source                 string
	UpdatedAt              sql.NullTime
}

func (q *Queries) CreateRecipe(ctx context.Context, arg CreateRecipeParams) (int64, error) {
	row := q.db.QueryRowContext(ctx, createRecipe,
		arg.Name,
		arg.Description,
		arg.InstructionsMarkdown,
		arg.ThumbnailUrl,
		arg.CookTimeSeconds,
		arg.PreparationTimeSeconds,
		arg.Source,
		arg.UpdatedAt,
	)
	var id int64
	err := row.Scan(&id)
	return id, err
}

const createRecipeIngredient = `-- name: CreateRecipeIngredient :exec
INSERT INTO recipe_ingredients (recipe_order, recipe_id, ingredient_id, unit, quantity)
VALUES (?, ?, ?, ?, ?)
`

type CreateRecipeIngredientParams struct {
	RecipeOrder  int64
	RecipeID     int64
	IngredientID int64
	Unit         string
	Quantity     int64
}

func (q *Queries) CreateRecipeIngredient(ctx context.Context, arg CreateRecipeIngredientParams) error {
	_, err := q.db.ExecContext(ctx, createRecipeIngredient,
		arg.RecipeOrder,
		arg.RecipeID,
		arg.IngredientID,
		arg.Unit,
		arg.Quantity,
	)
	return err
}

const createRecipeTag = `-- name: CreateRecipeTag :exec
INSERT INTO recipe_tags (recipe_order, recipe_id, tag_id)
VALUES (?, ?, ?)
`

type CreateRecipeTagParams struct {
	RecipeOrder int64
	RecipeID    int64
	TagID       int64
}

func (q *Queries) CreateRecipeTag(ctx context.Context, arg CreateRecipeTagParams) error {
	_, err := q.db.ExecContext(ctx, createRecipeTag, arg.RecipeOrder, arg.RecipeID, arg.TagID)
	return err
}

const createTag = `-- name: CreateTag :one
INSERT INTO tags (name)
VALUES (?)
RETURNING id
`

func (q *Queries) CreateTag(ctx context.Context, name string) (int64, error) {
	row := q.db.QueryRowContext(ctx, createTag, name)
	var id int64
	err := row.Scan(&id)
	return id, err
}

const getFullRecipe = `-- name: GetFullRecipe :one
SELECT id, name, description, instructions_markdown, thumbnail_url, servings, cook_time_seconds, preparation_time_seconds, source, created_at, updated_at, recipe_id, calories, fat, carbs, protein
FROM recipes
         JOIN main.recipe_nutrition rn on recipes.id = rn.recipe_id
WHERE recipes.id = ?
LIMIT 1
`

type GetFullRecipeRow struct {
	ID                     int64
	Name                   string
	Description            string
	InstructionsMarkdown   string
	ThumbnailUrl           sql.NullString
	Servings               int64
	CookTimeSeconds        int64
	PreparationTimeSeconds int64
	Source                 string
	CreatedAt              time.Time
	UpdatedAt              sql.NullTime
	RecipeID               int64
	Calories               float64
	Fat                    float64
	Carbs                  float64
	Protein                float64
}

func (q *Queries) GetFullRecipe(ctx context.Context, id int64) (GetFullRecipeRow, error) {
	row := q.db.QueryRowContext(ctx, getFullRecipe, id)
	var i GetFullRecipeRow
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Description,
		&i.InstructionsMarkdown,
		&i.ThumbnailUrl,
		&i.Servings,
		&i.CookTimeSeconds,
		&i.PreparationTimeSeconds,
		&i.Source,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.RecipeID,
		&i.Calories,
		&i.Fat,
		&i.Carbs,
		&i.Protein,
	)
	return i, err
}

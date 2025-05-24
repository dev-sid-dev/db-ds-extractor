//go:generate go run ./tooling/envs-gen "conf/app.toml" "conf/env.go" "conf"
//go:generate go fmt ./...

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/dev-sid-dev/db-ds-extractor/conf"
	"github.com/jackc/pgx/v5"
)

type Task struct {
	ID                      string
	CreatedAt               *time.Time
	UpdatedAt               *time.Time
	DeletedAt               *time.Time
	SKU                     *string
	MaterialID              *int64
	SceneID                 *int64
	LayoutID                *int64
	SurfaceID               *int64
	Lighting                *string
	Published               *bool
	Title                   *string
	ImageRendering          *string
	ImagePreview            *string
	SurfaceName             *string
	HasLayouts              *bool
	Lvl0                    *[]string
	Lvl1                    *[]string
	Lvl2                    *[]string
	Brand                   *string
	PrimaryColorFamilyLabel *string
	Product                 *json.RawMessage
	Price                   *string
}

func main() {
	sourceDB := conf.Get("SOURCE_DB")
	destDB := conf.Get("DESTINATION_DB")
	limitStr := conf.Get("LIMIT")

	fmt.Println("SOURCE_DB:", sourceDB)
	fmt.Println("DESTINATION_DB:", destDB)
	fmt.Println("LIMIT:", limitStr)

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		log.Fatalf("Invalid LIMIT value: %v", err)
	}

	srcConn, err := pgx.Connect(context.Background(), sourceDB)
	if err != nil {
		log.Fatalf("Failed to connect to source DB: %v", err)
	}
	defer srcConn.Close(context.Background())

	dstConn, err := pgx.Connect(context.Background(), destDB)
	if err != nil {
		log.Fatalf("Failed to connect to destination DB: %v", err)
	}
	defer dstConn.Close(context.Background())

	query := fmt.Sprintf(`
		SELECT id, created_at, updated_at, deleted_at, sku, material_id, scene_id, layout_id, surface_id,
		lighting, published, title, image_rendering, image_preview, surface_name, has_layouts, lvl0, lvl1, lvl2,
		brand, primary_color_family_label, product, price
		FROM public.tasks
		LIMIT %d
	`, limit)

	rows, err := srcConn.Query(context.Background(), query)
	if err != nil {
		log.Fatalf("Query error: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		t := Task{}
		err := rows.Scan(
			&t.ID, &t.CreatedAt, &t.UpdatedAt, &t.DeletedAt, &t.SKU,
			&t.MaterialID, &t.SceneID, &t.LayoutID, &t.SurfaceID,
			&t.Lighting, &t.Published, &t.Title, &t.ImageRendering,
			&t.ImagePreview, &t.SurfaceName, &t.HasLayouts,
			&t.Lvl0, &t.Lvl1, &t.Lvl2, &t.Brand,
			&t.PrimaryColorFamilyLabel, &t.Product, &t.Price,
		)
		if err != nil {
			log.Printf("Row scan error: %v", err)
			continue
		}

		_, err = dstConn.Exec(context.Background(), `
			INSERT INTO public.tasks (
				id, created_at, updated_at, deleted_at, sku, material_id, scene_id, layout_id, surface_id,
				lighting, published, title, image_rendering, image_preview, surface_name, has_layouts,
				lvl0, lvl1, lvl2, brand, primary_color_family_label, product, price
			)
			VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,
				$17,$18,$19,$20,$21,$22,$23)
		`,
			t.ID, t.CreatedAt, t.UpdatedAt, t.DeletedAt, t.SKU,
			t.MaterialID, t.SceneID, t.LayoutID, t.SurfaceID,
			t.Lighting, t.Published, t.Title, t.ImageRendering,
			t.ImagePreview, t.SurfaceName, t.HasLayouts,
			t.Lvl0, t.Lvl1, t.Lvl2, t.Brand,
			t.PrimaryColorFamilyLabel, t.Product, t.Price,
		)
		if err != nil {
			log.Printf("Insert error: %v", err)
		}
	}

	if err := rows.Err(); err != nil {
		log.Fatalf("Rows error: %v", err)
	}

	log.Println("Transfer completed successfully.")
}

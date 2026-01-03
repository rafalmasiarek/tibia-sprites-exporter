package app

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

// groupsIndex is a minimal JSON payload for downstream tooling.
type groupsIndex struct {
	GeneratedAt      string        `json:"generated_at"`
	AppearancesPath  string        `json:"appearances_path"`
	Groups           []groupsEntry  `json:"groups"`
	TotalGroups      int           `json:"total_groups"`
	NonEmptyGroups   int           `json:"non_empty_groups"`
}

type groupsEntry struct {
	FirstSpriteID int   `json:"first_sprite_id"`
	LastSpriteID  int   `json:"last_sprite_id"`
	Count         int   `json:"count"`
	SpriteIDs     []int `json:"sprite_ids"`
	OutputBase    string `json:"output_base"`
	OutputPNG     string `json:"output_png"`
}

// WriteGroupsIndexJSON writes a deterministic JSON file describing the groups extracted
// from the appearances file scan.
func WriteGroupsIndexJSON(outPath string, appearancesPath string, groups []spriteInfo) error {
	dir := filepath.Dir(outPath)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}

	payload := groupsIndex{
		GeneratedAt:     time.Now().UTC().Format(time.RFC3339),
		AppearancesPath: appearancesPath,
		Groups:          make([]groupsEntry, 0, len(groups)),
		TotalGroups:     len(groups),
	}

	nonEmpty := 0
	for _, g := range groups {
		if len(g.SpriteIDs) == 0 {
			continue
		}
		nonEmpty++

		first := g.SpriteIDs[0]
		last := g.SpriteIDs[len(g.SpriteIDs)-1]

		base := itoa(first)
		if first != last {
			base = itoa(first) + "-" + itoa(last)
		}

		payload.Groups = append(payload.Groups, groupsEntry{
			FirstSpriteID: first,
			LastSpriteID:  last,
			Count:         len(g.SpriteIDs),
			SpriteIDs:     g.SpriteIDs,
			OutputBase:    base,
			OutputPNG:     base + ".png",
		})
	}
	payload.NonEmptyGroups = nonEmpty

	b, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(outPath, b, 0o644)
}

func itoa(v int) string {
	// tiny helper to avoid importing strconv everywhere in this file
	return fmtInt(v)
}

func fmtInt(v int) string {
	if v == 0 {
		return "0"
	}
	neg := v < 0
	if neg {
		v = -v
	}
	buf := make([]byte, 0, 12)
	for v > 0 {
		d := v % 10
		buf = append(buf, byte('0'+d))
		v /= 10
	}
	if neg {
		buf = append(buf, '-')
	}
	// reverse
	for i, j := 0, len(buf)-1; i < j; i, j = i+1, j-1 {
		buf[i], buf[j] = buf[j], buf[i]
	}
	return string(buf)
}

package initfiles

import (
	"context"
	"encoding/json"
	"errors"
	"exp1/internal/types"
	"fmt"
	"os"
	"path"
)

func init(){
	InitFile(CreateConfig)
}

func CreateConfig(ctx context.Context, cancel context.CancelFunc) error {
	path := path.Join(".rec", "config.json")
    if ctx.Err() != nil {
        cancel()
        return ctx.Err()
    }

    // initial empty config
	repository := types.Repository{
		UserName: "",
		RemoteUrl: "",
	}

	recorder := types.Recorder{
		DebounceTime: 2,	//  initial default value
	}

    cfg := types.Config{
		Repository: repository,
		Recorder: recorder,
    }

    data, err := json.MarshalIndent(cfg, "", "  ")
    if err != nil {
        return fmt.Errorf("failed to marshal config: %w", err)
    }

    // create file with read and write permission for owner only (0644)
    f, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
    if err != nil {
        return fmt.Errorf("failed to create config file: %w", err)
    }
    defer f.Close()

    if _, err := f.Write(data); err != nil {
        return fmt.Errorf("failed to write config: %w", err)
    }

    // check cancellation AFTER writing
    if ctx.Err() != nil {
        cancel()
        return errors.New("operation canceled during config creation")
    }

    return nil
}
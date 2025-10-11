package recorder

import "math"

type Delta struct {
    SizeDiff int64
    HashDiff bool
}

func Compare(file *FileState) *Delta {
    return &Delta{
        SizeDiff: int64(math.Abs(float64(file.CurrentSize - file.PrevSize))),
        // HashDiff: file.Hash != file.Hash,
    }
}

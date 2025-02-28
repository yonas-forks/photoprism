package media

import (
	"strings"

	"github.com/jinzhu/gorm"

	"github.com/photoprism/photoprism/pkg/fs"
)

var PreviewFileTypes = []string{fs.ImageJpeg.String(), fs.ImagePng.String()}
var PreviewExpr = gorm.Expr("'" + strings.Join(PreviewFileTypes, "','") + "'")

package stores

import (
	"context"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"gitee.com/i-Things/share/conf"
	"gitee.com/i-Things/share/def"
	"github.com/twpayne/go-geom/encoding/ewkb"
	"github.com/twpayne/go-geom/encoding/ewkbhex"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"math"
)

type Point struct {
	Longitude float64 `json:"longitude,range=[0:180]"` //经度
	Latitude  float64 `json:"latitude,range=[0:90]"`   //纬度
}

// pgsql参考: https://www.jianshu.com/p/88ff6f693ffe?ivk_sa=1024320u
func (Point) GormDataType() string {
	switch dbType {
	case conf.Pgsql:
		return "GEOMETRY(point, 4326)"
	}
	return "point"
}

func (p Point) ToPo() def.Point {
	return def.Point{
		Longitude: p.Longitude,
		Latitude:  p.Latitude,
	}
}
func ToPoint(p def.Point) Point {
	return Point{
		Longitude: p.Longitude,
		Latitude:  p.Latitude,
	}
}
func (p *Point) Range(columnName string, Range int64) string {
	switch dbType {
	case conf.Pgsql:
		return fmt.Sprintf("ST_DWithin(%s,ST_GeomFromText('POINT(%v %v)', 4326),%v)",
			columnName, p.Longitude, p.Latitude, Range)
	default:
		return fmt.Sprintf(
			"round(st_distance_sphere(ST_GeomFromText('POINT(%v %v)'), ST_GeomFromText(AsText(%s))),2)>%d",
			p.Longitude, p.Latitude, columnName, Range)
	}
}

// hexToWKT converts a hex-encoded geometry to WKT format.
func hexToWKT(hexStr []byte) (string, error) {
	// Decode hex string to bytes.
	bytes, err := hex.DecodeString(string(hexStr))
	if err != nil {
		return "", err
	}

	// Try to parse the bytes as WKB (Well-Known Binary).
	g, err := ewkb.Unmarshal(bytes)
	if err != nil {
		return "", err
	}
	fmt.Println(g)
	// Convert the geometry to WKT (Well-Known Text).
	return "", nil
}

func (p *Point) parsePoint(binaryData []byte) error {
	if dbType == conf.Pgsql {
		g, err := ewkbhex.Decode(string(binaryData))
		if err != nil {
			return err
		}
		p.Longitude, p.Latitude = g.FlatCoords()[0], g.FlatCoords()[1]
		return nil
	}
	//下面是mysql的方式
	if len(binaryData) != 25 {
		return nil
	}
	longitudeBytes := binaryData[len(binaryData)-16 : len(binaryData)-8]
	latitudeBytes := binaryData[len(binaryData)-8:]
	var encode binary.ByteOrder = binary.LittleEndian
	if binaryData[4] != 1 {
		encode = binary.BigEndian
	}
	longitude := math.Float64frombits(encode.Uint64(longitudeBytes))
	latitude := math.Float64frombits(encode.Uint64(latitudeBytes))
	p.Longitude = longitude
	p.Latitude = latitude
	return nil
}
func (p *Point) Scan(value interface{}) error {
	if value == nil {
		return fmt.Errorf("failed to scan point: value is nil")
	}
	switch value.(type) {
	case []byte:
		va := value.([]byte)
		return p.parsePoint(va)
	case string:
		va := value.(string)

		return p.parsePoint([]byte(va))
	default:
		return fmt.Errorf("failed to scan point: invalid type: %T", value)
	}
	return nil
}

//func (p Point) Value() (driver.Value, error) {
//	return []byte(fmt.Sprintf("ST_GeomFromText('POINT(%f %f)')", p.Longitude, p.Latitude)), nil
//}

func (p Point) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
	switch dbType {
	case conf.Pgsql:
		return clause.Expr{
			//SQL:  "ST_PointFromText(?)",
			//
			SQL: fmt.Sprintf("ST_GeomFromText('POINT(%f %f)', 4326)", p.Longitude, p.Latitude), //如果你不知道 SRID 的值，可以使用 -1 来表示未知的空间参考系统。
		}
	default:
		return clause.Expr{
			SQL:  "ST_PointFromText(?)",
			Vars: []interface{}{fmt.Sprintf("POINT(%f %f)", p.Longitude, p.Latitude)},
		}
	}
}

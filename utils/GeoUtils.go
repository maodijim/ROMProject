package utils

import (
	Cmd "ROMProject/Cmds"
	log "github.com/sirupsen/logrus"
	"math"
)

const AtkRangeScale = 1000

func GetAngleByAxisY(src *Cmd.ScenePos, target *Cmd.ScenePos) float64 {
	return math.Atan2(float64(target.GetX()-src.GetX()), float64(target.GetZ()-src.GetZ())) * 57.29578
}

func CalcDir(angleY float64) float64 {
	dir := float64(int32(angleY)%360) + (angleY - float64(int32(angleY)))
	if dir < 0 {
		dir = dir + 360
	}
	return math.Floor(dir * AtkRangeScale)
}

func GetDistanceXZ(src *Cmd.ScenePos, target *Cmd.ScenePos) float64 {
	x := src.GetX() - target.GetX()
	z := src.GetZ() - target.GetZ()
	return math.Sqrt(float64(x*x + z*z))
}

func GetDistanceXZSquare(src *Cmd.ScenePos, target *Cmd.ScenePos) float64 {
	x := src.GetX() - target.GetX()
	z := src.GetZ() - target.GetZ()
	return float64(x*x + z*z)
}

func GetPosAwayFromTarget(src *Cmd.ScenePos, target *Cmd.ScenePos, disToTarget float64) *Cmd.ScenePos {
	angleY := CalcDir(GetAngleByAxisY(src, target))
	quadrant := GetQuadrantByAngle(angleY)
	angle := 0.0
	if quadrant == 1 || quadrant == 3 {
		angle = 90 - float64(int32(angleY)%90) + (angleY - float64(int32(angleY)))
	} else if quadrant == 2 || quadrant == 4 {
		angle = float64(int32(angleY)%90) + (angleY - float64(int32(angleY)))
	}

	targetPos := GetDistanceXZ(src, target)
	disFromX := int32(math.Ceil(math.Cos(DegreeToRadian(angle)) * (targetPos - disToTarget)))
	disFromY := int32(math.Ceil(math.Sin(DegreeToRadian(angle)) * (targetPos - disToTarget)))
	var x, y int32
	switch quadrant {
	case 1:
		x = src.GetX() + disFromX
		y = src.GetY() + disFromY
	case 2:
		x = src.GetX() + disFromX
		y = src.GetY() - disFromY
	case 3:
		x = src.GetX() - disFromX
		y = src.GetY() - disFromY
	case 4:
		x = src.GetX() - disFromX
		y = src.GetY() + disFromY
	default:
		log.Warnf("Invalid Quadrant %d angle %f angleY %f \n%v %v", quadrant, angle, angleY, src, target)
		return target
	}

	newPos := &Cmd.ScenePos{
		X: &x,
		Y: &y,
		Z: target.Z,
	}
	return newPos
}

func DegreeToRadian(degree float64) float64 {
	return degree * (math.Pi / 180)
}

func GetQuadrantByAngle(angle float64) int32 {
	angleInt := int(angle) % 360
	if angleInt >= 0 && angleInt <= 90 {
		return 1
	} else if angleInt < 360 && angleInt > 270 {
		return 2
	} else if angleInt <= 270 && angleInt > 180 {
		return 3
	} else if angleInt > 90 && angleInt <= 180 {
		return 4
	}
	return 0
}

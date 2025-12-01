package positionHelper

import (
	"context"
	"io"
	"os"
	"time"

	Cmd "ROMProject/Cmds"
	"ROMProject/gameConnection"
	"ROMProject/utils"

	log "github.com/sirupsen/logrus"
)

type PositionTask struct {
	GC        *gameConnection.GameConnection
	ctx       context.Context
	cancel    context.CancelFunc
	logWriter io.Writer
	logger    *log.Logger
}

func (p *PositionTask) GetContext() context.Context {
	return p.ctx
}

func (p *PositionTask) SetLogger(writer io.Writer) {
	mw := io.MultiWriter(p.GC.LogWriter(), writer)
	p.logWriter = mw
	p.logger.SetOutput(mw)
}

func (p *PositionTask) Start() {
	p.GC.ShouldChangeScene = true
	p.GC.GameServerLogin()

	go func() {
		for {
			select {
			case <-p.ctx.Done():
				p.logger.Info("停止自动跟随任务")
				p.cancel()
				p.GC.Close()
				return
			default:
				p.doFollow()
			}
		}
	}()

}

func (p *PositionTask) Stop() {
	p.cancel()
}

func (p *PositionTask) doFollow() {
	if p.GC.GetTeamLeaderName(true) != "" {
		leader := p.GC.GetTeamLeaderData(true)
		if utils.GetMemberDataByType(leader.GetDatas(), Cmd.EMemberData_EMEMBERDATA_MAPID) != uint64(p.GC.Role.GetMapId()) {
			p.logger.Infof("队长 %s (ID: %d) 不在当前地图", leader.GetName(), leader.GetGuid())
			p.GC.Role.FollowUserId = 0
		} else if p.GC.Role.FollowUserId != 0 {
			pos := p.GC.Role.GetPos()
			p.logger.Infof("正在跟随队长 %s (ID: %d) 坐标: x:%d y:%d z:%d", leader.GetName(), leader.GetGuid(), pos.GetX(), pos.GetY(), pos.GetZ())
		} else {
			p.logger.Infof("开始跟随队长 %s (ID: %d)", leader.GetName(), leader.GetGuid())
			p.GC.FollowUser(leader.GetGuid())
		}
	} else if p.GC.GetTeamLeaderName(true) == "" {
		p.logger.Infof("丢失队长，等待中...")
		p.GC.Role.FollowUserId = 0
	} else {
		p.logger.Infof("当前没有队长，等待中...")
	}
	time.Sleep(time.Second * 5)
}

func NewPositionTask(ctx context.Context, gc *gameConnection.GameConnection) *PositionTask {
	newCtx, cancel := context.WithCancel(ctx)
	mw := io.MultiWriter(os.Stdout, gc.LogWriter())
	logger := log.New()
	logger.SetOutput(mw)
	logger.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
	return &PositionTask{
		GC:     gc,
		ctx:    newCtx,
		cancel: cancel,
		logger: logger,
	}
}

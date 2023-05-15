package appeal

import (
	"fmt"
	"github.com/darabuchi/log"
	"github.com/darabuchi/utils"
	"github.com/darabuchi/utils/db"
	"github.com/darabuchi/utils/mq"
	"github.com/darabuchi/utils/xtime"
	"github.com/elliotchance/pie/v2"
	"github.com/xihui-forever/mutualRead/mail"
	"github.com/xihui-forever/mutualRead/types"
	"time"
)

func Load() error {
	mq.RegisterTopic(types.EventAppealChangedTopic, &types.EventAppealChanged{})

	mq.RegisterHandel(types.EventAppealChangedTopic, types.Topic("default"), &mq.Handle{
		HandleFunc: func(msg *mq.HandleReq) (*mq.HandleRsp, error) {
			var rsp mq.HandleRsp

			req := msg.Message.(*types.EventAppealChanged)
			a := req.Appeal

			switch a.State {
			case types.AppealStateWaitReviewer:
				var reviever types.ModelStudent
				err := db.Where("id = ?", a.ReviewerId).First(&reviever).Error
				if err != nil {
					log.Errorf("err:%v", err)
					return nil, err
				}

				err = mail.Send([]string{
					fmt.Sprintf("%s <%s>", reviever.Name, reviever.Email),
				}, "新申诉待处理通知", fmt.Sprintf("亲爱的 %s 同学，您好\n您有一条针对所阅试卷的申诉信息待处理 ,请于24小时之内处理完毕！\n后续系统会及时发送最终处理结果的邮件提醒，敬请关注！", reviever.Name))
				if err != nil {
					log.Errorf("err:%v", err)
					return nil, err
				}

			case types.AppealStateWaitTeacher:
				var teacher types.ModelTeacher
				err := db.Where("id = ?", a.TeacherId).First(&teacher).Error
				if err != nil {
					log.Errorf("err:%v", err)
					return nil, err
				}

				err = mail.Send([]string{
					fmt.Sprintf("%s <%s>", teacher.Name, teacher.Email),
				}, "新申诉待处理通知", fmt.Sprintf("尊敬的 %s 老师，您好！\n您录入的的试卷有一条相关申诉待处理，请登录系统查阅申诉的详细信息，并于一周内处理完毕，否则逾期申诉作废！\nITest团队", teacher.Name))
				if err != nil {
					log.Errorf("err:%v", err)
					return nil, err
				}

			case types.AppealStateFinish:
				var examiner types.ModelStudent
				err := db.Where("id = ?", a.ExaminerId).First(&examiner).Error
				if err != nil {
					log.Errorf("err:%v", err)
					return nil, err
				}

				err = mail.Send([]string{
					fmt.Sprintf("%s <%s>", examiner.Name, examiner.Email),
				}, "申诉处理结果通知", fmt.Sprintf("亲爱的 %s 同学，您好！\n您针对试卷提出的申诉最终处理结果已出，申诉结果为 %s，详细信息请登录系统进行查看！\nITest团队", examiner.Name, a.AppealResult))
				if err != nil {
					log.Errorf("err:%v", err)
					return nil, err
				}
			}

			return &rsp, nil
		},
		MaxAttempts:   5,
		MaxProcessCnt: 1,
		MsgTimeout:    time.Minute * 10,
	})

	go func() {
		err := Timeout()
		if err != nil {
			log.Errorf("err:%v", err)
		}

		for range time.NewTicker(xtime.Hour).C {
			err = Timeout()
			if err != nil {
				log.Errorf("err:%v", err)
			}
		}
	}()

	return nil
}

func Timeout() (err error) {
	// 24小时超时
	for {
		var appeals []*types.ModelAppeal
		err = db.Model(&types.ModelAppeal{}).
			Where("state = ?", types.AppealStateWaitReviewer).
			Where("created_at < ?", time.Now().Add(-time.Hour*24).Unix()).
			Limit(100).Find(&appeals).Error
		if err != nil {
			log.Errorf("err:%v", err)
			break
		}

		if len(appeals) == 0 {
			break
		}

		err = db.Model(&types.ModelAppeal{}).
			Where("state = ?", types.AppealStateWaitReviewer).
			Where("id in (?)", utils.PluckUint64(appeals, "Id")).
			Updates(map[string]interface{}{
				"state":       types.AppealStateTimeout,
				"reviewer_at": time.Now().Unix(),
			}).Error
		if err != nil {
			log.Errorf("err:%v", err)
			break
		}

		pie.Each(appeals, func(a *types.ModelAppeal) {
			a.State = types.AppealStateWaitTeacher
			_, err = mq.Publish(types.EventAppealChangedTopic, &types.EventAppealChanged{
				Appeal: a,
			})
			if err != nil {
				log.Errorf("err:%v", err)
			}
		})
	}

	for {
		var appeals []*types.ModelAppeal
		err = db.Model(&types.ModelAppeal{}).
			Where("state = ?", types.AppealStateWaitTeacher).
			Where("reviewer_at < ?", time.Now().Add(-xtime.Week).Unix()).
			Limit(100).Find(&appeals).Error
		if err != nil {
			log.Errorf("err:%v", err)
			break
		}

		if len(appeals) == 0 {
			break
		}

		err = db.Model(&types.ModelAppeal{}).
			Where("state = ?", types.AppealStateWaitTeacher).
			Where("id in (?)", utils.PluckUint64(appeals, "Id")).
			Updates(map[string]interface{}{
				"state":         types.AppealStateFinish,
				"reviewer_at":   time.Now().Unix(),
				"appeal_result": "超时未处理，申诉作废",
			}).Error
		if err != nil {
			log.Errorf("err:%v", err)
			break
		}

		pie.Each(appeals, func(a *types.ModelAppeal) {
			a.State = types.AppealStateFinish
			a.AppealResult = "超时未处理，申诉作废"

			_, err = mq.Publish(types.EventAppealChangedTopic, &types.EventAppealChanged{
				Appeal: a,
			})
			if err != nil {
				log.Errorf("err:%v", err)
			}
		})
	}

	return nil
}

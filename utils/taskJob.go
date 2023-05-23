package utils

import (
	"errors"
	"fmt"
	"ginDemo/global"
	"go.uber.org/zap"
)

// AddTaskJob 创建定时任务
func AddTaskJob(Job JobStruct) error {
	GlobalJobs = append(GlobalJobs, Job)
	_, err := global.JobS.Cron(Job.Expression).Do(func() {
		fmt.Printf("Running job %s (%s)\n", Job.Id, Job.Name)
	})
	if err != nil {
		zap.L().Error("新增Job任务失败", zap.Error(err))
		return err
	}
	return nil
}

// DelTaskJob 删除定时任务
func DelTaskJob(id string) error {
	for i, job := range GlobalJobs {
		if job.Id == id {
			GlobalJobs = append(GlobalJobs[:i], GlobalJobs[i+1:]...)
			err := global.JobS.RemoveByTag(fmt.Sprintf("job-%s", id))
			if err != nil {
				zap.L().Error("移除任务失败，请联系管理员!", zap.Error(err))
			}
		} else {
			return errors.New("暂未查询到该任务或该任务已调度完成！")
		}
	}
	return nil
}

// ChangeTaskJob 修改定时任务
func ChangeTaskJob(id string, Job JobStruct) error {
	for i, j := range GlobalJobs {
		if j.Id == id {
			GlobalJobs[i] = Job
			err := global.JobS.RemoveByTag(fmt.Sprintf("job-%s", id))
			if err != nil {
				zap.L().Error("移除任务失败，请联系管理员!", zap.Error(err))
				return errors.New("移除任务失败")
			}
			// 此处填写任务函数
			_, err = global.JobS.Cron(Job.Expression).Do(func() {
				fmt.Printf("Running job %d (%s)\n", Job.Id, Job.Name)
			})
			if err != nil {
				zap.L().Error("创建任务失败，请联系管理员!", zap.Error(err))
				return errors.New("移除任务成功，创建新任务失败！")
			}
		} else {
			return errors.New("暂未查询到该任务或该任务已调度完成！")
		}
	}
	return nil
}

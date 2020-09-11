package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/qwqcode/qwquiver/lib/exd"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// examCmd represents the db command
var examCmd = &cobra.Command{
	Use:     "exam [command]",
	Version: rootCmd.Version,
	Aliases: []string{},
	Short:   "考试管理",
	Long:    `qwquiver 的考试数据管理工具`,
	Args:    cobra.NoArgs,
}

var examListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "列出所有 Exam",
	Run:     examListRun,
	Args:    cobra.NoArgs,
}

var examRemoveCmd = &cobra.Command{
	Use:     "remove [ExamName]",
	Aliases: []string{"rm", "remove", "del"},
	Short:   "删除指定的 Exam",
	Run:     examRemoveRun,
	Args:    cobra.RangeArgs(1, 1),
}

// $ qwquiver exam config
var examConfigCmd = &cobra.Command{
	Use:     "config [Command]",
	Aliases: []string{"conf"},
	Short:   "配置指定的 Exam",
	Args:    cobra.NoArgs,
}

// $ qwquiver exam config set
var configSetCmd = &cobra.Command{
	Use:   "set [ExamName]",
	Short: "配置指定的 Exam",
	Args:  cobra.RangeArgs(1, 2),
	Run:   examConfSetRun,
}

// $ qwquiver exam config get
var configGetCmd = &cobra.Command{
	Use:   "get [ExamName]",
	Short: "获取指定 Exam 的配置",
	Args:  cobra.RangeArgs(1, 1),
	Run:   examConfigGetRun,
}

func init() {
	rootCmd.AddCommand(examCmd)

	// list
	examCmd.AddCommand(examListCmd)

	// remove
	examCmd.AddCommand(examRemoveCmd)
	flagP(examRemoveCmd, "force", "f", false, "强制删除，没有任何提示")

	// config
	examCmd.AddCommand(examConfigCmd)

	// config set
	examConfigCmd.AddCommand(configSetCmd)
	flag(configSetCmd, "Grp", "", "修改 Exam Grp")
	flag(configSetCmd, "Label", "", "修改 Exam Label")
	flag(configSetCmd, "Subj", "", "修改 Exam Subj (JSON格式)")
	flag(configSetCmd, "SubjFullScore", "", "修改 Exam SubjFullScore (JSON格式)")
	flag(configSetCmd, "Date", "", "修改 Exam Date")
	flag(configSetCmd, "Note", "", "修改 Exam Note")
	flag(configSetCmd, "json", "", "直接通过 JSON 格式修改所有参数")

	// config get
	examConfigCmd.AddCommand(configGetCmd)
	flagP(configGetCmd, "inline", "i", false, "输出未格式化的 JSON 数据")
}

func examListRun(cmd *cobra.Command, args []string) {
	allExams := exd.GetAllExams()
	fmt.Printf("共有 %d 个 Exam\n", len(allExams))
	for i, exam := range allExams {
		itemLen := exam.CountScores()
		fmt.Printf(" %d. %s (%d 条数据)\n", i+1, exam.Name, itemLen)
	}
}

func examRemoveRun(cmd *cobra.Command, args []string) {
	examName := args[0]

	if force, _ := cmd.Flags().GetBool("force"); !force {
		fmt.Println("请加上 flag '--force' 确认删除 '" + examName + "' 这个 Exam")
		return
	}

	err := exd.RemoveExam(examName)
	if err != nil {
		fmt.Println("删除 Exam 发生错误")
		fmt.Println(err)
	}
}

// config set
func examConfSetRun(cmd *cobra.Command, args []string) {
	examName := args[0]
	if !exd.HasExam(examName) {
		logrus.Error("Exam '" + examName + "' 不存在")
		return
	}

	examConf := exd.GetExam(examName)

	jsonStr, _ := cmd.Flags().GetString("json")
	if jsonStr == "" {
		jsonStr = args[1]
	}

	if jsonStr != "" {
		if err := exd.UpdateExamConfByJSON(examConf, jsonStr); err != nil {
			logrus.Error("保存 ExamConf 发生错误 ", err)
			return
		}
	} else {
		// 非直接 JSON 修改模式
		if grp, _ := cmd.Flags().GetString("Grp"); grp != "" {
			examConf.Grp = grp
		}
		if label, _ := cmd.Flags().GetString("Label"); label != "" {
			examConf.Label = label
		}
		if subjStr, _ := cmd.Flags().GetString("Subj"); subjStr != "" {
			var subjList []string
			err := json.Unmarshal([]byte(subjStr), &subjList)
			if err == nil && subjList != nil {
				examConf.Subj = subjStr
			} else {
				logrus.Error("尝试解析 flag '--Subj' 的 JSON 数据时发生错误 ", err)
			}
		}
		if sfc, _ := cmd.Flags().GetString("SubjFullScore"); sfc != "" {
			var subjFullScore map[string]float64
			err := json.Unmarshal([]byte(sfc), &subjFullScore)
			if err == nil && subjFullScore != nil {
				examConf.SubjFullScore = sfc
			} else {
				logrus.Error("尝试解析 flag '--SubjFullScore' 的 JSON 数据时发生错误 ", err)
			}
		}
		if date, _ := cmd.Flags().GetString("Date"); date != "" {
			examConf.Date = date
		}
		if note, _ := cmd.Flags().GetString("Note"); note != "" {
			examConf.Note = note
		}
		// TODO ... 比较 dirty 的代码，先将就这样
		exd.UpdateExamConf(examConf)
	}

	examConfJSON := exd.GetExamConfJSONStr(examName, true)
	fmt.Println("ExamConf 已更新")
	fmt.Println(string(examConfJSON))
}

// config get
func examConfigGetRun(cmd *cobra.Command, args []string) {
	examName := args[0]
	if !exd.HasExam(examName) {
		logrus.Error("Exam '" + examName + "' 不存在")
		return
	}

	examConfJSON := "NULL"
	if isInlineJSON, _ := cmd.Flags().GetBool("inline"); isInlineJSON {
		examConfJSON = exd.GetExamConfJSONStr(examName, false)
	} else {
		examConfJSON = exd.GetExamConfJSONStr(examName, true)
	}

	fmt.Println(string(examConfJSON))
}

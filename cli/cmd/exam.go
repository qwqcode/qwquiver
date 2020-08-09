package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/qwqcode/qwquiver/lib"
	"github.com/qwqcode/qwquiver/model"
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

func init() {
	rootCmd.AddCommand(examCmd)

	var examListCmd = &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "列出所有 Exam",
		Run:     cmdExamListRun,
		Args:    cobra.NoArgs,
	}
	examCmd.AddCommand(examListCmd)

	var examRemoveCmd = &cobra.Command{
		Use:     "remove [ExamName]",
		Aliases: []string{"rm", "remove", "del"},
		Short:   "删除指定的 Exam",
		Run:     cmdExamRemoveRun,
		Args:    cobra.RangeArgs(1, 1),
	}
	flagP(examRemoveCmd, "force", "f", false, "强制删除，没有任何提示")
	examCmd.AddCommand(examRemoveCmd)

	// $ qwquiver exam config
	var examConfigCmd = &cobra.Command{
		Use:     "config [Command]",
		Aliases: []string{"conf"},
		Short:   "配置指定的 Exam",
		Args:    cobra.NoArgs,
	}
	examCmd.AddCommand(examConfigCmd)

	// $ qwquiver exam config set
	var configSetCmd = &cobra.Command{
		Use:   "set [ExamName]",
		Short: "配置指定的 Exam",
		Args:  cobra.RangeArgs(1, 1),
		Run:   cmdExamConfigSetRun,
	}
	flag(configSetCmd, "grp", "", "修改 Exam Grp")
	flag(configSetCmd, "label", "", "修改 Exam Label")
	flag(configSetCmd, "subj", "", "修改 Exam Subj (JSON格式)")
	flag(configSetCmd, "subj-full-score", "", "修改 Exam SubjFullScore (JSON格式)")
	flag(configSetCmd, "date", "", "修改 Exam Date")
	flag(configSetCmd, "note", "", "修改 Exam Note")
	flag(configSetCmd, "json", "", "直接通过 JSON 格式修改所有参数")
	examConfigCmd.AddCommand(configSetCmd)

	var configGetCmd = &cobra.Command{
		Use:   "get [ExamName]",
		Short: "获取指定 Exam 的配置",
		Args:  cobra.RangeArgs(1, 1),
		Run:   cmdExamConfigGetRun,
	}
	flagP(configGetCmd, "inline", "i", false, "输出未格式化的 JSON 数据")
	examConfigCmd.AddCommand(configGetCmd)
}

func cmdExamListRun(cmd *cobra.Command, args []string) {
	allExamNames := lib.GetAllExamNames()
	fmt.Printf("共有 %d 个 Exam\n", len(allExamNames))
	for i, name := range allExamNames {
		itemLen := lib.GetExamScoreLen(name)
		fmt.Printf(" %d. %s (%d 条数据)\n", i+1, name, itemLen)
	}
}

func cmdExamRemoveRun(cmd *cobra.Command, args []string) {
	examName := args[0]

	if force, _ := cmd.Flags().GetBool("force"); !force {
		fmt.Println("请加上 flag '--force' 确认删除 '" + examName + "' 这个 Exam")
		return
	}

	err := lib.RemoveExam(examName)
	if err != nil {
		fmt.Println("删除 Exam 发生错误")
		fmt.Println(err)
	}
}

// config set
func cmdExamConfigSetRun(cmd *cobra.Command, args []string) {
	examName := args[0]
	if !lib.IsExamExist(examName) {
		logrus.Error("Exam '" + examName + "' 不存在")
		return
	}

	examConf := lib.GetExamConf(examName)
	if examConf == nil {
		// examConf 不存则创建一个新的
		examConf = &model.ExamConf{Name: examName}
		if err := lib.SaveExamConf(examConf); err != nil {
			logrus.Error("创建 ExamConf 发生错误 ", err)
			return
		}
	}

	if jsonStr, _ := cmd.Flags().GetString("json"); jsonStr != "" {
		if err := lib.UpdateExamConfByJSON(examConf, jsonStr); err != nil {
			logrus.Error("保存 ExamConf 发生错误 ", err)
			return
		}
	} else {
		// 非直接 JSON 修改模式
		if grp, _ := cmd.Flags().GetString("grp"); grp != "" {
			examConf.Grp = grp
		}
		if label, _ := cmd.Flags().GetString("label"); label != "" {
			examConf.Label = label
		}
		if subjStr, _ := cmd.Flags().GetString("subj"); subjStr != "" {
			var subjList []string
			err := json.Unmarshal([]byte(subjStr), &subjList)
			if err == nil && subjList != nil {
				examConf.Subj = subjList
			} else {
				logrus.Error("尝试解析 flag '--subj' 的 JSON 数据时发生错误 ", err)
			}
		}
		if sfc, _ := cmd.Flags().GetString("subj-full-score"); sfc != "" {
			var subjFullScore map[string]float64
			err := json.Unmarshal([]byte(sfc), &subjFullScore)
			if err == nil && subjFullScore != nil {
				examConf.SubjFullScore = subjFullScore
			} else {
				logrus.Error("尝试解析 flag '--subj-full-score' 的 JSON 数据时发生错误 ", err)
			}
		}
		if date, _ := cmd.Flags().GetString("date"); date != "" {
			examConf.Date = date
		}
		if note, _ := cmd.Flags().GetString("note"); note != "" {
			examConf.Note = note
		}
		// TODO ... 比较 dirty 的代码，先将就这样
		lib.SaveExamConf(examConf)
	}

	examConfJSON := lib.GetExamConfJSONStr(examName, true)
	fmt.Println("ExamConf 已更新")
	fmt.Println(string(examConfJSON))
}

// config get
func cmdExamConfigGetRun(cmd *cobra.Command, args []string) {
	examName := args[0]
	if !lib.IsExamExist(examName) {
		logrus.Error("Exam '" + examName + "' 不存在")
		return
	}

	examConfJSON := "NULL"
	if isInlineJSON, _ := cmd.Flags().GetBool("inline"); isInlineJSON {
		examConfJSON = lib.GetExamConfJSONStr(examName, false)
	} else {
		examConfJSON = lib.GetExamConfJSONStr(examName, true)
	}

	fmt.Println(string(examConfJSON))
}

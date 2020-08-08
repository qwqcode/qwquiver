package cmd

import (
	"fmt"

	"github.com/qwqcode/qwquiver/lib"
	"github.com/spf13/cobra"
)

// dbCmd represents the db command
var dbCmd = &cobra.Command{
	Use:     "db [command]",
	Version: rootCmd.Version,
	Aliases: []string{"database"},
	Short:   "数据库管理",
	Long:    `qwquiver 的 bbolt 数据管理工具`,
	Args:    cobra.NoArgs,
}

func init() {
	rootCmd.AddCommand(dbCmd)
	flagP(dbCmd, "name", "n", "", "数据表名")

	var dbListCmd = &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "列出所有 bucket",
		Run:     cmdDbListRun,
		Args:    cobra.NoArgs,
	}
	dbCmd.AddCommand(dbListCmd)

	var dbRemoveCmd = &cobra.Command{
		Use:     "remove [BucketName]",
		Aliases: []string{"rm", "remove", "del"},
		Short:   "删除指定的 bucket",
		Run:     cmdDbRemoveRun,
		Args:    cobra.RangeArgs(1, 1),
	}
	flagP(dbRemoveCmd, "force", "f", false, "强制删除，没有任何提示")
	dbCmd.AddCommand(dbRemoveCmd)
}

func cmdDbListRun(cmd *cobra.Command, args []string) {
	allBucketNames := lib.GetAllBucketNames()
	fmt.Printf("共有 %d 个 Bucket\n", len(allBucketNames))
	for i, name := range allBucketNames {
		itemLenS := ""
		if itemLen := lib.GetExamScoreLen(name); itemLen > 0 {
			itemLenS = fmt.Sprintf("(%d 条 Score 数据)", itemLen)
		}
		fmt.Println(fmt.Sprintf(" %d. %s %s", i+1, name, itemLenS))
	}
}

func cmdDbRemoveRun(cmd *cobra.Command, args []string) {
	bucketName := args[0]

	if force, _ := cmd.Flags().GetBool("force"); !force {
		fmt.Println("请加上 flag '--force' 确认删除 '" + bucketName + "' 这个 Bucket")
		return
	}

	err := lib.RemoveBucket(bucketName)
	if err != nil {
		fmt.Println("删除 Bucket 发生错误")
		fmt.Println(err)
	}
}

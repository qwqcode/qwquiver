<?php

namespace app\commands;

use app\app\ScoreHandle;
use yii\console\Controller;
use yii\helpers\Console;

class ScoreController extends Controller
{
    /**
     * 根据成绩构建排名
     *
     * @param string $tableId 数据表 ID
     * @return int
     */
    public function actionRanking($tableId = 'ALL')
    {
        return $this->_actionTpl($tableId, function (ScoreHandle $score) use (&$tableId) {
            $errors = $score->makeRanking();
            if (empty($errors)) {
                echo "数据表 score_{$tableId} 的排名已生成\n\n";
            }
            else {
                echo "数据表 score_{$tableId} 的排名生成失败:\n";
                var_dump($errors);
                echo "\n";
            }
        });
    }
    
    /**
     * 根据排名构ID（唯一的）
     *
     * @param string $fieldName 数据生成目标字段名（字段类型 bigint，不能填 自动递增 的字段）
     * @param string $tableId 数据表 ID
     * @return int
     */
    public function actionSortId($fieldName, $tableId = 'ALL')
    {
        return $this->_actionTpl($tableId, function (ScoreHandle $score) use (&$fieldName, &$tableId) {
            $errors = $score->makeSortId($fieldName);
            if (empty($errors)) {
                echo "score_{$tableId} 排序ID 已生成\n\n";
            }
            else {
                echo "score_{$tableId} 排序ID 生成失败:\n";
                var_dump($errors);
                echo "\n";
            }
            
        }, function () use (&$fieldName) {
            if (empty($fieldName)) {
                echo '$fieldName 不能为空';
                
                return false;
            }
            
            return true;
        });
    }
    
    /**
     * 统一指定表原本不同的学校名和班级名
     *
     * @param string $targetTableId 需要被规范的表 ID
     * @throws \Exception
     */
    public function actionJoinableSCName($targetTableId)
    {
        if ($targetTableId == 'ALL') {
            foreach (ScoreHandle::map() as $id => $item)
                ScoreHandle::makeJoinableSchoolClassName($id, true);
        } else {
            ScoreHandle::makeJoinableSchoolClassName($targetTableId);
        }
        
        Console::output('学校名和班级名格式统一完毕');
    }
    
    /**
     * 操作模板
     *
     * @param string $tableId
     * @param callable $action
     * @param callable|null $beforeDo
     * @param callable|null $afterDo
     * @return int
     */
    public function _actionTpl($tableId = 'ALL', callable $action, callable $beforeDo = null, callable $afterDo = null)
    {
        if (!is_null($beforeDo) && !$beforeDo())
            return 1;
        
        $do = function ($tableId) use (&$action, &$times, &$total) {
            if (!ScoreHandle::isExist($tableId)) {
                echo '表 ID=' . $tableId . ' 不存在';
                
                return;
            }
    
            echo '[INFO] 开始处理表 ID=' . $tableId . ' 的数据';
            $score = new ScoreHandle($tableId);
            $action($score);
        };
        
        if ($tableId == 'ALL') {
            foreach (ScoreHandle::map() as $id => $item)
                $do($id);
        }
        else {
            $do($tableId);
        }
        
        if (!is_null($afterDo) && !$afterDo())
            return 1;
        
        return 0;
    }
}

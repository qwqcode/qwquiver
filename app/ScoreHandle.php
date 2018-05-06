<?php
namespace app\app;


use yii\db\Command;
use yii\db\Query;
use yii\helpers\Console;

class ScoreHandle extends Score
{
    /**
     * 构建排名
     *
     * 成绩  ID B
     * 691  1  1
     * 691  2  1
     * 548  3  3
     * 444  4  4
     * 444  5  4 [important]
     * 333  6  6
     *
     * @return array 执行错误列表
     */
    public function makeRanking()
    {
        $error = [];
        $findAll = $this->find()
            ->orderBy(['overall'=>SORT_DESC])
            ->all();
        
        $ranking = 0;
        $rankingOffset = 1;
        $overallTemp = null;
        
        $total = count($findAll);
        $done = 0;
        Console::startProgress($done, $total);
        foreach ($findAll as $item) {
            if (count($error) > 10)
                break;
            if ($item['overall']!==$overallTemp) {
                $ranking = $ranking + $rankingOffset;
                $rankingOffset = 1;
            } else {
                $rankingOffset ++;
            }
            
            try {
                \Yii::$app->db->createCommand()
                    ->update($this->_formName, ['ranking' => $ranking], ['id'=>$item['id']])
                    ->execute();
                
            } catch (\Exception $exception) {
                $error[$item['id']] = $exception->getMessage();
            }
            
            $overallTemp = $item['overall'];
            
            $done++;
            Console::updateProgress($done, $total);
        }
        Console::endProgress();
        
        return $error;
    }
    
    /**
     * 排序规则构建ID
     *
     * @param $fieldName string 数据修改字段名
     * @param $orderBy array
     * @return array 执行错误列表
     */
    public function makeSortId($fieldName, $orderBy=['overall'=>SORT_DESC])
    {
        $error = [];
        $findAll = $this->find()
            ->orderBy($orderBy)
            ->all();
        
        $id = 0;
        
        $total = count($findAll);
        $done = 0;
        Console::startProgress($done, $total);
        foreach ($findAll as $item) {
            if (count($error) > 10)
                break;
            
            $id ++;
            
            try {
                \Yii::$app->db->createCommand()
                    ->update($this->_formName, ["$fieldName" => $id], ["id"=>$item['id']])
                    ->execute();
                
            } catch (\Exception $exception) {
                $error[$item['id']] = $exception->getMessage();
            }
            
            $done++;
            Console::updateProgress($done, $total);
        }
        Console::endProgress();
        
        return $error;
    }
    
    /**
     * 统一原本不同的学校名和班级名
     *
     * @param string $targetTableId 需要被处理的数据表
     * @param bool $force 强制，无提示
     * @throws \Exception
     */
    public static function makeJoinableSchoolClassName($targetTableId, $force=false)
    {
        define('NormalTableId', self::getParams('normalTableId'));
        if ($targetTableId === NormalTableId) {
            Console::error('目标表ID 不能等于 模范表ID');
            return;
        }
        
        $schoolDictionary = [];
        $classDictionary = [];
        $dictionary = [];
        
        // 调出规范表中独一无二的姓名
        $normalTableName = 'score_' . NormalTableId;
        $targetTableName = 'score_' . $targetTableId;
        
        $normalTableItems = (new Query())
            ->select([XM['FN'], XX['FN'], BJ['FN']])
            ->from($normalTableName)
            ->orderBy([XX['FN'] => SORT_ASC, BJ['FN'] => SORT_ASC])
            ->all();
        
        // 如果目标表中找不到这个，直接删除
        $targetTableItems = (new Query())
            ->select([XM['FN'], XX['FN'], BJ['FN']])
            ->from($targetTableName)
            ->orderBy([XX['FN'] => SORT_ASC, BJ['FN'] => SORT_ASC])
            ->all();
        
        $getNormalUniqueItems = function ($items) use ($targetTableItems) {
            $arr_raw = $items;
            $arr = $arr_raw;
            
            $foreachHandle = function (&$arr, callable $do) {
                foreach ($arr as $k => $v) {
                    $name = $v[XM['FN']];
                    $school = $v[XX['FN']];
                    $class = $v[BJ['FN']];
                    $do($k, $v, $name, $school, $class);
                }
            };
            
            $names = [];
            $foreachHandle($arr, function ($k, $v, $name, $school, $class) use (&$names, &$arr, $targetTableItems) {
                // 如果名字是 2 个字，直接删除
                if (mb_strlen($name, 'utf8') <= 2) {
                    unset($arr[$k]);
                    return;
                }
                
                // 如果目标表里没有，或非 独一无二 的名字，直接删除
                $targetTableFindCount = count(array_filter($targetTableItems, function ($item) use ($name) {
                    return ($item[XM['FN']] == $name);
                }));
                if ($targetTableFindCount <= 0 || $targetTableFindCount > 1) {
                    unset($arr[$k]);
                    return;
                }
                
                // 如果模范表里没有，或非 独一无二 的名字，直接删除
                $normalTableFindCount = count(array_filter($arr, function ($item) use ($name) {
                    return ($item[XM['FN']] == $name);
                }));
                if ($normalTableFindCount <= 0 || $normalTableFindCount > 1) {
                    unset($arr[$k]);
                    return;
                }
            });
            
            // 班级+学校 保留一个相同
            $t2 = [];
            $foreachHandle($arr, function ($k, $v, $name, $school, $class) use (&$t2, &$arr) {
                if (!in_array($school . $class, $t2)) {
                    $t2[] = $school . $class;
                } else {
                    unset($arr[$k]);
                }
            });
            
            return $arr;
        };
        
        // 成功获取到最佳模范样品
        $normalUniqueItems = $getNormalUniqueItems($normalTableItems);
        //var_dump($normalUniqueItems);die();
        
        // 调出目标表中数据，并建立词典
        foreach ($normalUniqueItems as $index => $item) {
            $targetItem = (new Query())
                ->select([XM['FN'], XX['FN'], BJ['FN']])
                ->from($targetTableName)
                ->where([XM['FN'] => $item[XM['FN']]])
                ->one();
            
            $schoolDictionary[$targetItem[XX['FN']]] = $item[XX['FN']];
            $classDictionary[$targetItem[XX['FN']]][$targetItem[BJ['FN']]] = $item[BJ['FN']];
            
            $dictionary[] = [
                $targetItem[XX['FN']],
                $targetItem[BJ['FN']],
                $item[XX['FN']],
                $item[BJ['FN']]
            ];
        }
        
        //var_dump($schoolDictionary);
        //var_dump($classDictionary);
        
        $updateQueries = [];
        
        foreach ($dictionary as $item) {
            $beforeSchool = $item[0];
            $beforeClass = $item[1];
            $afterSchool = $item[2];
            $afterClass = $item[3];
            
            $query = (new Query())
                ->createCommand()
                ->update($targetTableName, [
                    XX['FN'] => $afterSchool,
                    BJ['FN'] => $afterClass
                ], [
                    XX['FN'] => $beforeSchool,
                    BJ['FN'] => $beforeClass
                ]);
            
            Console::output($query->getRawSql());
            $updateQueries[] = $query;
        }
        
        Console::output();
        
        if (!$force) {
            if (!Console::confirm('确定要执行这些 SQL语句 吗？')) {
                Console::error('已取消执行 SQL语句');
                
                return;
            }
        }
        
        $total = count($dictionary);
        $done = 0;
        Console::startProgress($done, $total);
        foreach ($updateQueries as $query) {
            /** @var $query Command */
            $query->execute();
            $done++;
            Console::updateProgress($done, $total);
        }
        Console::endProgress();
    }
}
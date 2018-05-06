<?php
namespace app\app;

use yii\base\InvalidValueException;
use yii\db\Query;
use yii\helpers\Console;

class Score
{
    protected $_id = null; // ID
    protected $_formName = null; // 完整数据表名
    
    /**
     * Score 构造函数
     * @param $fId
     */
    function __construct($fId)
    {
        if (!$this->isExist($fId)) {
            throw new InvalidValueException('$formName must be in the formMap()');
        }
        
        $this->_id = $fId;
        $this->_formName = 'score_' . $fId;
    }
    
    /**
     * 获取名称
     *
     * @return string
     */
    public function getLabel()
    {
        return self::map()[$this->_id]['label'];
    }
    
    /**
     * 索引
     *
     * @return array
     */
    public static function map()
    {
        $map = self::getParams('map');
        
        $mapHandle = [];
        foreach ($map as $key=>$value) {
            $mapHandle[$key]['label'] = $value['label'];
            $mapHandle[$key]['fields'] = array_merge([BH, XM, XX, BJ, ZF, PM], $value['fields']);
        }
        
        return $mapHandle;
    }
    
    /**
     * 获取参数
     *
     * @param null $key
     * @return null
     */
    public static function getParams($key=null)
    {
        if (!is_null($key)) {
            return \Yii::$app->params['score'][$key] ?? null;
        } else {
            return \Yii::$app->params['score'];
        }
    }
    
    /**
     * id是否在map中存在
     *
     * @param $id
     * @return bool
     */
    public static function isExist($id)
    {
        return isset(self::map()[$id]) ? true : false;
    }
    
    /**
     * 获取 GET 请求 ID 参数
     *
     * @param null $defaultValue
     * @return array|mixed
     */
    public static function getId($defaultValue = null)
    {
        return \Yii::$app->request->get(self::getIdReqParName(), $defaultValue);
    }
    
    /**
     * Get 请求 ID 参数名
     *
     * @return string
     */
    public static function getIdReqParName()
    {
        return 'id';
    }
    
    /**
     * 获取默认 ID
     *
     * @return mixed
     */
    public static function getDefaultId()
    {
        foreach (self::map() as $key => $value)
            return $key;
        
        return null;
    }
    
    /**
     * 获取 “数据列表” 项目，用于列表显示
     *
     * @return array
     */
    public static function getHtmlUlItems()
    {
        $items = [];
        
        foreach (self::map() as $name => $item) {
            $items[] = ['label'=>$item['label'], 'url'=>['site/index', self::getIdReqParName()=>$name]];
        }
        
        return $items;
    }
    
    /**
     * 获取数据的所有字段
     *
     * @param bool $onlyIntField 仅仅返回 int 格式的字段
     * @return array
     */
    public function getFields($onlyIntField=false)
    {
        $fields = self::map()[$this->_id]['fields'];
        
        if ($onlyIntField) {
            $intFields = [];
            foreach ($fields as $field) {
                if (in_array($field, INT_FIELDS))
                    $intFields[] = $field;
            }
            return $intFields;
        }
        
        return $fields;
    }
    
    /**
     * 字段名是否存在
     *
     * @param $FN string 字段名
     * @return bool
     */
    public function isFieldNameExist($FN)
    {
        $result = false;
        foreach ($this->getFields() as $item) {
            if ($item['FN'] == $FN) {
                $result = true;
                break;
            }
        }
        
        return $result;
    }
    
    /**
     * Get Query
     *
     * @return Query|null
     */
    public function find()
    {
        return (new Query())
            ->select('*')
            ->from($this->_formName);
    }
    
    /**
     * 查找
     *
     * @param array $where 筛查规则
     * @param int $page 页码数
     * @param array $orderBy 排序规则
     * @param int $pagePer 每页项目
     * @return array
     */
    public function queryGetApi($where, $page, $orderBy=[ZF['FN']=>SORT_DESC], $pagePer = 50)
    {
        $offset = ($page == 1) ? 0 : ($page - 1) * $pagePer;
        
        $query = self::fastQuery($where, $pagePer, $offset, $orderBy)->all();
        
        $lastPageNum = $this->getLastPageNum($where, $pagePer);
        $data = self::api($query, $pagePer, $page, $lastPageNum, $orderBy);
        
        return array_merge(['dataTitle' => $this->getLabel()], $data);
    }
    
    /**
     * 快速查询
     *
     * @param $where
     * @param $limit
     * @param $offset
     * @param $orderBy
     * @return Query|null
     */
    public function fastQuery($where, $limit=null, $offset=null, $orderBy=null)
    {
        $find = $this->find();
        
        if (!is_null($where))
            $find = $find->filterWhere(array_merge($where, []));
        
        if (!is_null($limit))
            $find = $find->limit($limit);
        
        if (!is_null($offset))
            $find = $find->offset($offset);
        
        if (!is_null($orderBy))
            $find = $find->orderBy($orderBy);
        
        return $find;
    }
    
    /**
     * 获取最后一页的页码
     *
     * @param $where
     * @param $pagePer
     * @return int
     */
    public function getLastPageNum($where, $pagePer)
    {
        $query = $this->fastQuery($where, null, null, null);
        
        return intval(ceil(intval($query -> count()) / intval($pagePer)));
    }
    
    /**
     * 生成Api数据
     *
     * @param $rawObj array 源数据
     * @param int $pagePer 每页数量
     * @param int $nowPage 当前页码
     * @param int $lastPage 最后一页页码
     * @param array $orderBy 排序规则
     * @return array
     */
    public function api($rawObj, $pagePer, $nowPage, $lastPage, $orderBy)
    {
        $apiArr = [];
        $total = 0;
        $inArrFieldName = $this->getFields();
        if ($nowPage<=$lastPage && $nowPage>=1) {
            $num = 0;
            foreach ($rawObj as $key => $score) {
                $theDataNum = (($nowPage - 1) * $pagePer) + ($key + 1);
                // 基本信息
                $apiArr[$num][BH['FN']] = $theDataNum;
                foreach ($inArrFieldName as $fieldNameItem) {
                    
                    // 排除编号
                    if ($fieldNameItem == BH) {
                        continue;
                    }
                    
                    $apiArr[$num][$fieldNameItem['FN']] = $score[$fieldNameItem['FN']];
                    
                    // DEMO
                    $len = mb_strlen($score[$fieldNameItem['FN']]);
                    $testStr = '';
                    for ($i=0; $i<$len; $i++)
                        $testStr .= 'X';
                    
                    $apiArr[$num][$fieldNameItem['FN']] = $testStr;
                }
                $num ++;
            }
            $total = count($apiArr);
        }
        
        $pagination = [
            'pagePer' => $pagePer,
            'nextBtnEnabled' => ($total >= intval($pagePer)) ? true : false, // 是否还有下一页
            'nowPage' => intval($nowPage),
            'lastPage' => $lastPage,
        ];
        
        return [
            'label' => $this->getLabel(),
            'total' => count($apiArr),
            'fieldList' => $inArrFieldName,
            'pagination' => $pagination,
            'sortBy' => $orderBy,
            'score' => $apiArr,
        ];
    }
    
    /**
     * 获取平均分
     *
     * @param $query Query
     * @param null|int $round 小数点后保留多少位
     * @return array
     */
    public function getAvgs($query, $round=null)
    {
        $cacheKey = md5($query->createCommand()->getRawSql());
        
        $getAvgByFieldName = function ($fieldName) use ($query, $cacheKey, $round)
        {
            $result = \Yii::$app->cache->getOrSet("avg_cache_{$cacheKey}_{$fieldName}", function() use ($query, $fieldName){
                return $query->average($fieldName);
            });
            
            if (!is_null($round)) {
                return round($result, $round);
            } else {
                return $result;
            }
        };
        
        $avgs = [];
        foreach ($this->getFields(true) as $item) {
            $fieldName = $item['FN'];
            $avgs[$fieldName] = [$getAvgByFieldName($fieldName), $item];
        }
        
        return $avgs;
    }
    
    /**
     * 获取分数
     *
     * @param $query Query
     *
     * @return array
     */
    public function getScores($query)
    {
        $scores = [];
        foreach ($this->getFields(true) as $item) {
            $fieldName = $item['FN'];
            $scores[$fieldName] = [intval($query->one()[$fieldName]), $item];
        }
        
        return $scores;
    }
}
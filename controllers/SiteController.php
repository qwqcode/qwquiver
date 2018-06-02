<?php
namespace app\controllers;

use app\app\Score;
use PhpOffice\PhpSpreadsheet\Spreadsheet;
use PhpOffice\PhpSpreadsheet\Writer\Xls;
use Yii;
use yii\helpers\Html;
use yii\web\Controller;
use yii\web\NotFoundHttpException;
use yii\web\Response;

class SiteController extends Controller
{
    /**
     * @inheritdoc
     */
    public function actions()
    {
        return [
            'error' => [
                'class' => 'yii\web\ErrorAction',
            ],
        ];
    }
    
    /**
     * 总览
     *
     * @throws \Exception
     * @return mixed
     */
    public function actionIndex()
    {
        $id = Score::getId(Score::getDefaultId()); // 数据表ID
        
        $page = intval(Yii::$app->request->get('page', 1)); // 页码
        $pagePer = Yii::$app->request->get('pagePer', 50); // 每页显示
        $sortBy = Yii::$app->request->get('sortBy'); // 排序规则
        $sortType = intval(Yii::$app->request->get('sortType', SORT_DESC)); // 排序类型 SORT_ASC|SORT_DESC
        $queryType = Yii::$app->request->get('queryType'); // 查询类型
        $studentName = Yii::$app->request->get('queryData');
        $schoolName = Yii::$app->request->get('queryData_school');
        $className = Yii::$app->request->get('queryData_class');
        $isFileSave = Yii::$app->request->get('saveFile') == 'y'; // 文件保存模式
        $fileSaveMode = $isFileSave ? Yii::$app->request->get('saveMode', 'paging') : 'paging'; // 保存模式
        
        // 若数据表不存在
        if (!Score::isExist($id))
            throw new NotFoundHttpException();
        
        // 如果不是数字
        if (!is_numeric($pagePer))
            throw new NotFoundHttpException();
        
        // 实例化 Score 对象
        $score = new Score($id);
        
        // 数据筛查
        if ($queryType === 'class') {
            $where = [XX['FN'] => $schoolName, BJ['FN'] => $className];
            $dataSubtitle = "{$schoolName} {$className} · 班级成绩";
        } else if ($queryType === 'school') {
            $where = [XX['FN'] => $schoolName];
            $dataSubtitle = "{$schoolName} · 全校成绩";
        } else if ($queryType === 'name') {
            $where = ['like', XM['FN'], $studentName];
            $dataSubtitle = "姓名包含 “{$studentName}” 的学生成绩";
        } else {
            $where = [];
            $dataSubtitle = '全市考生成绩';
        }
        
        // 排序规则
        if (!empty($sortBy) && $score->isFieldNameExist($sortBy)) {
            $sortType = in_array($sortType, [SORT_ASC, SORT_DESC]) ? $sortType : SORT_DESC;
            $orderBy = [$sortBy=>$sortType];
        } else {
            $orderBy = [ZF['FN']=>SORT_DESC];
        }
    
        // 文件保存模式，无分页
        if ($fileSaveMode == 'noPaging')
            $pagePer = 10000000;
        
        $data = array_merge($score->queryGetApi($where, $page, $orderBy, $pagePer), []);
        
        $data['dataSubtitle'] = $dataSubtitle;
        $data['dataSubtitleB'] = '[页码 ' . $page . '/' . $data['pagination']['lastPage'] . ']';
        $data['dataAvg'] = $score->getAvgs($score->fastQuery($where));
    
        // 文件保存模式
        if ($isFileSave) {
            return $this->saveExcelFile($data);
        }
        
        // Ajax 请求
        if (Yii::$app->request->isAjax) {
            Yii::$app->response->format = Response::FORMAT_JSON;
            Yii::$app->response->headers->set('Pragma', 'no-cache');
            Yii::$app->response->headers->set('Cache-Control', 'no-cache');
            
            return ['success'=>true, 'data'=>$data];
        }
        
        // 若没有任何成绩数据
        if (empty($data['score'])) {
            throw new NotFoundHttpException();
        }
        
        return $this->render('index', [
            'id' => $id,
            'score' => $score,
            'dataLabel' => Html::encode($score->getLabel()),
            'data' => $data,
        ]);
    }
    
    /**
     * 另存为 Excel 文件
     *
     * @param $_data
     * @throws \PhpOffice\PhpSpreadsheet\Exception
     * @throws \PhpOffice\PhpSpreadsheet\Writer\Exception
     * @return mixed
     */
    private function saveExcelFile($_data)
    {
        $fieldList = $_data['fieldList'];
        
        // 创建新的表格
        $spreadsheet = new Spreadsheet();
        
        // 用于 int 转 字母
        $a2z = [];
        for ($i = ord('A'); $i < ord('Z'); $i++)
            $a2z[] = chr($i);
        
        $rowTotal = 0;
        $setContent = function ($x, $y, $value) use ($spreadsheet, $a2z, &$rowTotal) {
            $spreadsheet
                ->setActiveSheetIndex(0)
                ->setCellValue($a2z[$x] . '' . ($y + 1), $value);
            
            $rowTotal++;
        };
        
        // 表头
        foreach ($fieldList as $x => $item) {
            $setContent($x, 0, $item['LB']);
        }
        
        // 数据
        foreach ($_data['score'] as $index => $item) {
            foreach ($fieldList as $x => $field)
                $setContent($x, $index + 1, $item[$field['FN']]);
        }
        
        // 设置宽度
        $spreadsheet->getActiveSheet()->getColumnDimension('A')->setWidth(8);
        $spreadsheet->getActiveSheet()->getColumnDimension('B')->setWidth(8);
        $spreadsheet->getActiveSheet()->getColumnDimension('C')->setWidth(20);
        $spreadsheet->getActiveSheet()->getColumnDimension('D')->setWidth(10);
        $spreadsheet->getActiveSheet()->setAutoFilter('C1:D' . $rowTotal);
        
        $writer = new Xls($spreadsheet);
        $filename = "{$_data['dataTitle']}（{$_data['dataSubtitle']}）.xls";
        header('Content-Type: application/vnd.ms-excel');
        header('Content-Disposition: attachment;filename="' . $filename . '.xls"'); /*-- $filename is  xsl filename ---*/
        header('Cache-Control: max-age=0');
        
        ob_end_clean();
        $writer->save('php://output');
        
        return '';
    }
    
    /**
     * 获取学校/班级所有分类
     *
     * @return mixed
     * @throws NotFoundHttpException
     */
    public function actionGetAllCategory()
    {
        if (Yii::$app->request->isAjax) {
            Yii::$app->response->format = Response::FORMAT_JSON;
            Yii::$app->response->headers->set('Pragma', 'no-cache');
            Yii::$app->response->headers->set('Cache-Control', 'no-cache');
        } else {
            throw new NotFoundHttpException();
        }
        
        $id = Score::getId(Score::getDefaultId()); // 数据表ID
        $school = trim(Yii::$app->request->get('school')); // 指定学校，为了获得所有班级分类
        
        // 若数据表不存在
        if (!Score::isExist($id))
            throw new NotFoundHttpException();
        
        // 实例化 Score 对象
        $score = new Score($id);
        
        if (empty($school)) {
            $groupBy = $score->find()->groupBy('school')->all();
            $results = [];
            foreach ($groupBy as $item) {
                if (in_array($item, $results))
                    continue;
                
                $results[] = $item['school'];
            }
        } else {
            $groupBy = $score->find()->where(['school' => $school])->groupBy('clas')->all();
            $results = [];
            foreach ($groupBy as $item) {
                if (in_array($item, $results))
                    continue;
                
                $results[] = $item['clas'];
            }
        }
        
        return ['success'=>true, 'data'=>$results];
    }
    
    /**
     * 图表
     */
    public function actionCharts()
    {
        $queryName = Yii::$app->request->get('queryName');
        $querySchool = Yii::$app->request->get('querySchool');
        $queryClass = Yii::$app->request->get('queryClass');
        $onlyFiledName = Yii::$app->request->get('onlyFiledName');
        
        $scoreMap = Score::map();
        $scoreAvg = [];
        $find = null;
        $subtitle = '';
        
        if (empty($queryName)) {
            $dataDesc = '近年市平均趋势';
        } else {
            $dataDesc = "{$queryName} 的" . ($onlyFiledName ?? "") . "成绩趋势";
            $subtitle = "{$querySchool} {$queryClass}";
        }
        
        foreach ($scoreMap as $key=>$item) {
            $score = new Score($key);
            
            $arr = [
                'name' => $score->getLabel(),
            ];
            
            if (empty($queryName)) {
                $find = $score->find();
                $avgs = $score->getAvgs($find, 2);
                foreach ($avgs as $avgKey=>$avgItem) {
                    $arr[$avgItem[1]['LB']] =  $avgItem[0];
                };
            } else {
                $find = $score->find()->where([
                    XM['FN'] => $queryName,
                    XX['FN'] => $querySchool,
                    BJ['FN'] => $queryClass,
                ]);
                if ($find->count() > 1)
                    $subtitle .= '，数据可能是假的，因为名字不是独一无二的';
                
                $data = $score->getScores($find);
                foreach ($data as $dataKey=>$dataItem) {
                    $arr[$dataItem[1]['LB']] =  $dataItem[0];
                };
                $arr[PM['LB']] = -intval($find->one()[PM['FN']]);
            }
            
            $scoreAvg[] = $arr;
        }
        
        return $this->render('charts', [
            'dataDesc' => $dataDesc,
            'subtitle' => $subtitle,
            'scoreAvg' => array_reverse($scoreAvg),
            'onlyFiledName' => $onlyFiledName
        ]);
    }
    
    /**
     * 关于
     *
     * @return string
     */
    public function actionAbout()
    {
        return $this->render('about');
    }
}

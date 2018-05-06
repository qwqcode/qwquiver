<?php

/* @var $this yii\web\View */
/* @var $id string */
/* @var $score app\app\Score */
/* @var $dataLabel string */
/* @var $data array */

use yii\helpers\Json;
use app\app\Score;

$this->title = $dataLabel;
?>
<div class="site-index content-inner">
    
    <div class="grades-table card">
        <div class="card-header">
            <h2 class="card-title" data-wlytable="title"></h2>
            <small class="card-subtitle">本次考试共有 <?= $score->find()->count() ?> 人参加</small>
            <div class="actions">
                <span onclick="wlySearch.showPanel()" class="actions__item show-top-badge">
                    <i class="zmdi zmdi-search"></i> <span>搜索</span>
                </span>
                <span onclick="wlyTableDataSave.show()" class="actions__item">
                    <i class="zmdi zmdi-download"></i> <span>下载</span>
                </span>
                <span data-wly-toggle="wlyTablePrint" class="actions__item">
                    <i class="zmdi zmdi-print"></i> <span>打印</span>
                </span>
                <span onclick="wlyTable.DisplayController.show()" class="actions__item">
                    <i class="zmdi zmdi-format-paint"></i> <span>表格调整</span>
                </span>
                <span onclick="window.wlyTable.showDataCounter()" class="actions__item">
                    <i class="zmdi zmdi-flash"></i> <span>平均分</span>
                </span>
                <span data-wly-toggle="wlyTableFullScreen" class="actions__item">
                    <i class="zmdi zmdi-fullscreen"></i> <span>全屏显示</span>
                </span>
            </div>
        </div>
        <div class="card-block" style="padding: 0;">
            <div class="wly-table-container" data-toggle="wlyTable" style="opacity: 0;height: 500px"></div>
        </div>
    </div>
    
</div>

<script>
<?php $this->beginBlock('script'); ?>
(function () {
    window.wlyTableDataId = '<?= $id ?>'; // 数据表 ID
    window.wlyTableConfig = <?= Json::encode([
        'data' => $data,
        'dataIdList' => Score::getHtmlUlItems(),
        'sign' => \Yii::$app->params['sign'],
    ], JSON_FORCE_OBJECT); ?>;
    
    $(document).ready(function(){
        wlyTable.init();
        // 加载完毕
        wlyPageLoader();
    });
})();
<?php $this->endBlock(); ?>
</script>
<?php $this->registerJs($this->blocks['script'], $this::POS_HEAD); ?>


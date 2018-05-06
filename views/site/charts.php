<?php

/* @var $this yii\web\View */
/* @var $dataDesc string */
/* @var $scoreAvg array */
/* @var $subtitle string */
/* @var $onlyFiledName array|mixed */

use yii\helpers\Json;
use yii\helpers\Html;

$this->title = $dataDesc . ' - 图表趋势';
/*$this->params['breadcrumbs'][] = $this->title;*/
?>
<div class="site-charts">
    <div class="card">
        <div class="card-header">
            <h2 class="card-title"><?= Html::encode($dataDesc) ?><?php if (!$subtitle): ?><small class="card-subtitle"><?= Html::encode($subtitle) ?></small><?php endif; ?></h2>
            
            <div class="actions">
                <span onclick="toggleShowTotalScore(this);" class="actions__item">
                    <i class="zmdi zmdi-square-o"></i> 仅显示总成绩
                </span>
                <span onclick="toggleOnlyShowFinalScore(this);" class="actions__item">
                    <i class="zmdi zmdi-square-o"></i> 仅显示毕业月考成绩
                </span>
            </div>
        </div>
        <div class="card-block">
            <div class="chart-wrap"></div>
        </div>
    </div>
</div>

<script>
    var rawData = <?= Json::encode($scoreAvg); ?>;
    var _onlyFiledName = null;
    <?php if (!empty($onlyFiledName)): ?>
    _onlyFiledName = <?= Json::encode([$onlyFiledName ?? '']); ?>;
    <?php endif; ?>
    
    $(document).ready(function () {
        buildChart(rawData, _onlyFiledName);
    });
  
    function getFinalScoreData() {
    var data = [];
    for (var i in rawData) {
        if (rawData[i]['name'].indexOf('毕业') > -1) {
            data.push(rawData[i]);
        }
    }
    return data;
  }
  
    var isShowTotalScore = false;
    var isOnlyShowFinalScore = false;
    
    function toggleShowTotalScore(btnObj) {
        if (isShowTotalScore) {
            $(btnObj).find('i').attr('class', 'zmdi zmdi-square-o');
            isShowTotalScore = false;
        } else {
            $(btnObj).find('i').attr('class', 'zmdi zmdi-check-square');
            isShowTotalScore = true;
        }
      
        rebuildChart();
    }

    function toggleOnlyShowFinalScore(btnObj) {
        if (isOnlyShowFinalScore) {
            isOnlyShowFinalScore = false;
            $(btnObj).find('i').attr('class', 'zmdi zmdi-square-o');
        } else {
            isOnlyShowFinalScore = true;
            $(btnObj).find('i').attr('class', 'zmdi zmdi-check-square');
        }
      
        rebuildChart();
    }
  
    function rebuildChart() {
        var data = rawData;
        var fields = null;
        if (isShowTotalScore)
            fields = ['总分'];
        if (isOnlyShowFinalScore)
            data = getFinalScoreData();
            
        buildChart(data, fields);
    }
  
    function buildChart(data, fields) {
        
        if (typeof(fields) === "undefined" || !fields) {
            fields = ['语文', '数学', '英语', '物理', '化学', '政治', '历史', '地理', '生物'];
        }

        window._chart = null;
        $('.chart-wrap').html('<div id="myChart"></div>');
        
        var ds = new DataSet();
        var dv = ds.createView().source(data);
        dv.transform({
            type: 'fold',
            fields: fields, // 展开字段集
            key: 'subject', // key字段
            value: 'score', // value字段
        });
  
        var chart = new G2.Chart({
            container: 'myChart',
            forceFit: true,
            height : getHeight()
        });
        chart.source(dv, {
            name: {
                range: [ 0, 1 ]
            }
        });
        chart.tooltip({
            crosshairs: {
                type: 'line'
            }
        });
        chart.axis('score', {
            label: {
                formatter: val => {
                    return val + ' 分';
                }
            }
        });
        chart.line().position('name*score').color('subject').shape('circle');
        chart.point().position('name*score').color('subject').size(4).shape('circle').style({
            stroke: '#fff',
            lineWidth: 1
        });
        chart.render();

        window._chart = chart;
    }
    
    $(window).resize(function () {
        chartFitScreen();
    });

    window.wlySidebar.onAfterToggle.push(function () {
        chartFitScreen();
    });
    
    function chartFitScreen() {
        window._chart.forceFit();
        window._chart.changeHeight(getHeight());
    }
</script>
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
                <span style="margin-right: 20px" id="show_fields">
                    <span onclick="toggleShowFields(this, ['总分']);" class="actions__item">
                        <i class="zmdi zmdi-square-o"></i> 仅总分
                    </span>
                    <span onclick="toggleShowFields(this, ['市排名']);" class="actions__item">
                        <i class="zmdi zmdi-square-o"></i> 仅排名
                    </span>
                </span>
                <span onclick="toggleOnlyShowFinalScore(this);" class="actions__item">
                    <i class="zmdi zmdi-square-o"></i> 仅月考
                </span>
            </div>
        </div>
        <div class="card-block">
            <div class="chart-wrap"></div>
        </div>
    </div>
</div>

<script>
    var _rawData = <?= Json::encode($scoreAvg); ?>;
    
    var _showFiledNames = null;
    var _showFiledNamesRestDefault = function () {
        _showFiledNames = ['语文', '数学', '英语', '物理', '化学', '政治', '历史', '地理', '生物'];
    };
    
    <?php if (empty($onlyFiledName)): ?>
        _showFiledNamesRestDefault();
    <?php else: ?>
        _showFiledNames = <?= Json::encode([$onlyFiledName ?? '']); ?>;
    <?php endif; ?>
    
    $(document).ready(function () {
        buildChart(_rawData, _showFiledNames);
    });

    function toggleShowFields(btnObj, fieldNames) {
        var on = btnCheck(btnObj);
        if (!on) {
            _showFiledNames = fieldNames;
        } else {
            _showFiledNamesRestDefault();
        }
        rebuildChart();
        $('#show_fields').find('span').each(function (i, item) {
            btnCheck(item, false);
        });
        btnCheck(btnObj, !on);
    }
    
    var isOnlyShowFinalScore = false;
    function toggleOnlyShowFinalScore(btnObj) {
        var on = btnCheck(btnObj);
        isOnlyShowFinalScore = !on;
        rebuildChart();
        btnCheck(btnObj, !on);
    }
    
    function btnCheck(btnObj, needCheck) {
        if (typeof needCheck === 'boolean') {
            if (needCheck) {
                $(btnObj).find('i').attr('class', 'zmdi zmdi-check-square');
            } else {
                $(btnObj).find('i').attr('class', 'zmdi zmdi-square-o');
            }
        } else {
            return $(btnObj).find('i').is('.zmdi.zmdi-check-square');
        }
    }
  
    function rebuildChart() {
        var data = _rawData;
        
        if (isOnlyShowFinalScore) {
            data = (function getFinalScoreData() {
                var arr = [];
                for (var i in _rawData) {
                    if (_rawData[i]['name'].indexOf('毕业') > -1) {
                        arr.push(_rawData[i]);
                    }
                }
                return arr;
            })();
        }
        
        buildChart(data, _showFiledNames);
    }
  
    function buildChart(data, fields) {
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
                formatter: function (val) {
                    return val;
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
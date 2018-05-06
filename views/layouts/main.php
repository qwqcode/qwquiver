<?php

/* @var $this \yii\web\View */
/* @var $content string */

use app\app\ToolBox;
use app\app\Score;
use yii\helpers\Html;
use yii\helpers\Url;
use app\app\MainAsset;

MainAsset::register($this);

$atHome = $this->context->id == 'site' && $this->context->action->id == 'index';
$baseActions = [
    ['icon'=>'view-carousel', 'label'=>'总览', 'url'=>['site/index']],
    ['icon'=>'equalizer', 'label'=>'趋势', 'url'=>['site/charts']],
    ['icon'=>'info-outline', 'label'=>'关于', 'url'=>['site/about']]
];
?>
<?php $this->beginPage() ?>
<!DOCTYPE html>
<html lang="<?= Yii::$app->language ?>">
<head>
    <!-- { "QWQUERY", "https://github.com/Zneiat/qwquery" } -->
    <meta charset="<?= Yii::$app->charset ?>">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width,minimum-scale=1.0">
    <meta name="copyright" content="qwqaq.com">
    <meta name="referrer" content="always">
    <meta name="renderer" content="webkit">
    <link href="/image/favicon.png" rel="shortcut icon" type="image/x-icon">
    <?= Html::csrfMetaTags() ?>
    <title><?= Html::encode($this->title) ?> - <?= Yii::$app->params['title'] ?></title>
    <?php if (!YII_DEBUG): ?>
    <script>
    var _hmt = _hmt || [];
    (function() {
        var hm = document.createElement("script");
        hm.src = "https://hm.baidu.com/hm.js?edf0f175c0cb7427ec9a2a712ba95228";
        var s = document.getElementsByTagName("script")[0];
        s.parentNode.insertBefore(hm, s);
    })();
    </script>
    <?php endif; ?>
    <?php $this->head() ?>
</head>
<body>
<?php $this->beginBody() ?>

<div class="top-header">
    <div class="main-navbar">
        <div class="left">
            <div class="sidebar-toggle-btn"><i class="zmdi zmdi-menu"></i></div>
            <span class="brand"><a href="/" class="app-name"><?= Yii::$app->params['title'] ?></a></span>
        </div>
        <div class="right">
            <ul class="actions-btn-bar main-bar">
                <?php if ($atHome): ?>
                <li><a onclick="wlySearch.showPanel()"><i class="zmdi zmdi-search"></i> <span>搜索</span></a></li>
                <?php endif; ?>
            </ul>
            <ul class="actions-btn-bar scene-bar" style="display: none;"></ul>
        </div>
    </div>
</div>

<div class="wrap">
    <div class="sidebar">
        <div class="widget link-list">
            <h2 class="list-label">基本</h2>
            <?= Html::ul($baseActions, ['item' => function($item, $index) {
                $options['class'] = ToolBox::isAtThisUrl($item['url']) ? 'active' : null;
                return Html::tag('li', Html::a("<i class=\"zmdi zmdi-{$item['icon']}\"></i> ".$item['label'], Url::toRoute($item['url'])), $options);
            }]);
            ?>
            <h2 class="list-label">数据列表</h2>
            <?= Html::ul(Score::getHtmlUlItems(), ['item' => function($item, $index) use ($atHome) {
                $options['class'] = (
                    (ToolBox::isAtThisUrl($item['url'])) ||
                    ($atHome && !Score::getId() && (Score::getDefaultId() == $item['url'][Score::getIdReqParName()]))
                ) ? 'active' : null;
                return Html::tag('li', Html::a('<i class="zmdi zmdi-trending-up"></i> '.$item['label'], Url::toRoute($item['url'])), $options);
            }]);
            ?>
        </div>
    </div>
    
    <div class="main-content-area">
        <?= $content ?>
    </div>
</div>

<?php $this->endBody() ?>
</body>
</html>
<?php $this->endPage() ?>

<?php

/* @var $this yii\web\View */
/* @var $name string */
/* @var $message string */
/* @var $exception yii\web\NotFoundHttpException */

use yii\helpers\Html;

$this->title = $exception->statusCode . " " . $exception->getName();

if ($exception->statusCode == 404 && empty($message))
    $message = Yii::t('yii', 'Page not found.');
?>
<div class="site-error">

    <div class="card">
        <div class="card-block">
            <h2 class="card-title"><?= Html::encode($this->title) ?> - <?= nl2br(Html::encode($message)) ?></h2>
            <!--<br/>
            <p>Web服务器处理您的请求时发生上述错误。</p>
            <p>如果您认为这是服务器错误，请与我们联系。 谢谢。</p>-->
        </div>
    </div>
    
</div>
<script>
    $(document).ready(function () {
        wlyNavbar.actionsSceneBarSet();
        setTimeout(function () {
            wlyNavbar.sceneChange('<?= Html::encode($this->title) ?>', '#ff6b68');
        }, 200);
    });
</script>
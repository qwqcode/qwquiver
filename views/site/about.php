<?php

/* @var $this yii\web\View */

use yii\helpers\Html;

$this->title = '关于';
?>
<div class="site-about">
    <div class="card">
        <div class="card-header">
            <h2 class="card-title"><?= Html::encode($this->title) ?></h2>
        </div>
        <div class="card-block">
            <h4><a href="https://github.com/Zneiat/qwquery" target="_blank">QwQuery</a> (Version: <?= \Yii::$app->version ?>)</h4>
            <br/>
            <p>网站由 <a href="https://github.com/Zneiat/qwquery" target="_blank">QwQuery</a> 强力驱动，作者：</p>
            <ul>
                <li>Blog: <a href="https://qwqaq.com" target="_blank">QwQAQ.com</a></li>
                <li>GitHub: <a href="https://github.com/Zneiat" target="_blank">ZNEIAT</a></li>
                <li>E-mail: 1149527164@qq.com</li>
            </ul>
            <br/>
            <p>未经允许代码和衍生品不得用于商业用途，侵权必究</p>
            <br/>
            <p><?= Yii::$app->params['sign'] ?></p>
            <p><?= Yii::$app->params['statement'] ?></p>
        </div>
    </div>
</div>
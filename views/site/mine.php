<?php

/* @var $this yii\web\View */

$this->title = '个人分析';
?>
<div class="site-mine">
    <div class="card">
        <div class="card-block">
            <h2 style="text-align: center">时间不够 以后再做这个功能... <img src="/image/face_tuxie.jpg" style="display: block;margin: 35px auto 0 auto;"></h2>
        </div>
    </div>
</div>
<script>
    $(document).ready(function () {
        wlyNavbar.actionsSceneBarSet();
        setTimeout(function () {
            wlyNavbar.sceneChange('Wisely Query Version 3', '#f5c942');
        }, 200);
    });
</script>
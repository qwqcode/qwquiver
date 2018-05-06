<?php
namespace app\app;

use yii\bootstrap\BootstrapAsset;
use yii\bootstrap\BootstrapPluginAsset;
use yii\web\AssetBundle;
use yii\web\View;
use yii\web\YiiAsset;

/**
 * Main Asset
 */
class MainAsset extends AssetBundle
{
    public $basePath = '@webroot';
    public $baseUrl = '@web';

    public $css = [
        'css/qwquery.css',
        'css/material-design-iconic-font.min.css'
    ];

    public $js = [
        'js/qwquery.js',
        'js/jquery.history.js',
        'js/jQuery.print.min.js',
        'js/data-set.min.js',
        'js/g2.min.js'
    ];

    public $jsOptions = [
        'position' => View::POS_HEAD
    ];
    
    public $depends = [
        YiiAsset::class,
        BootstrapAsset::class,
        BootstrapPluginAsset::class,
    ];
}

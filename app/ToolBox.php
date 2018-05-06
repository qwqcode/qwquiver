<?php
namespace app\app;

use Yii;

class ToolBox
{
    /**
     * 判断当页是否为指定链接
     * For Checks whether a menu item is active.
     * @param $url
     * @return bool
     */
    public static function isAtThisUrl($url)
    {
        $_route = Yii::$app->controller->getRoute();
        $_params = Yii::$app->request->getQueryParams();
        
        if (isset($url) && is_array($url) && isset($url[0]))
        {
            $route = $url[0];
            if ($route[0] !== '/' && Yii::$app->controller)
            {
                $route = Yii::$app->controller->module->getUniqueId() . '/' . $route;
            }
            if (ltrim($route, '/') !== $_route)
            {
                return false;
            }
            unset($url['#']);
            if (count($url) > 1)
            {
                $params = $url;
                unset($params[0]);
                foreach ($params as $name => $value)
                {
                    if ($value !== null && (!isset($_params[$name]) || $_params[$name] != $value))
                    {
                        return false;
                    }
                }
            }
            
            return true;
        }
        
        return false;
    }
    
    /**
     * 字符转 UTF-8 格式
     *
     * @param $data
     * @return mixed|string
     */
    public static function charsetUtf8($data){
        if(!empty($data))
        {
            $fileType = mb_detect_encoding($data , array('UTF-8','GBK','LATIN1','BIG5')) ;
            if( $fileType != 'UTF-8'){
                $data = mb_convert_encoding($data ,'utf-8' , $fileType);
            }
        }
        return $data;
    }
    
    public static function getCsvFile($csvfile, $lines) {
        if(!$fp = fopen($csvfile, 'r')) {
            return false;
        }
        $i = 0;
        $data = [];
        while(($i++ < $lines) && !feof($fp)) {
            $data[] = fgetcsv($fp);
        }
        fclose($fp);
        return $data;
    }
}
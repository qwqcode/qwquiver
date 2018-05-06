<?php
$map = [];
$map['93_2_201804'] = [
    'label' => '毕业第二次月考',
    'fields' => [YW, SX, YY, WL, HX, ZZ, LS]
];
$map['93_1_201803'] = [
    'label' => '毕业第一次月考',
    'fields' => [YW, SX, YY, WL, HX, ZZ, LS]
];

return [
	'title' => '2018届初中成绩统计',
    'sign' => 'QwQAQ.com',
    'statement' => '数据均来源于网络，仅供参考',
    'score' => [
        'map' => $map,
        'normalTableId' => '93_1_201803',
    ],
    'cookieValidationKey' => '',
];

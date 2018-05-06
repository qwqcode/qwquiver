/**
 * QWQUERY
 * @author ZNEIAT<1149527164@qq.com>
 */

(function () {
    "use strict";

    var SORT_ASC = 4;
    var SORT_DESC = 3;

    window.getHeight = function () {
        return $(window).height() - $(wlyNavbar.$container).outerHeight(true) - $('.card .card-header').outerHeight(true) - 80;
    };

    /**
     * 初始化
     */
    $(document).ready(function () {
        // 当 URL 发生改变时
        History.Adapter.bind(window, 'statechange', function () {
            wlyTable.reqData(getRequest(), false);
        });

        wlyWrap.init();
        wlyNavbar.init();
        wlySidebar.init();
    });

    // 加载指示器
    window.wlyPageLoader = function () {
        var sLoader = "[data-toggle='page-loader']",
            $loader = $(sLoader);

        if ($loader.css('display') === 'none') {
            $loader.fadeIn()
        } else {
            $loader.fadeOut();
        }
    };

    // Wrap
    window.wlyWrap = {
        sWrap: '.wrap',
        sMainContentArea: '.main-content-area',

        init: function () {
            this.$mainContentArea = $(this.sMainContentArea);
        },
        changeMainContentAreaFull: function (needFull) {
            if (typeof (needFull) === 'boolean') {
                if (needFull)
                    wlyWrap.$mainContentArea.addClass('main-content-area-full');
                else
                    wlyWrap.$mainContentArea.removeClass('main-content-area-full');
            }
        }
    };

    // 导航栏
    window.wlyNavbar = {
        sContainer: '.main-navbar',
        sAppName: '.main-navbar .app-name',
        sActionsBarMain: '.main-navbar .actions-btn-bar.main-bar',
        sActionsBarScene: '.main-navbar .actions-btn-bar.scene-bar',
        appNameValue: '',
        bgValue: '',
        isShow: true,

        init: function () {
            this.$container = $(this.sContainer);
            this.$appName = $(this.sAppName);
            this.$actionsBarMain = $(this.sActionsBarMain);
            this.$actionsBarScene = $(this.sActionsBarScene);

            this.appNameValue = this.$appName.html();
            this.bgValue = this.$container.css('background');
        },
        // AppName 设置
        appNameSet: function (value) {
            this._appNameChange(value);
        },
        // AppName 复原
        appNameRec: function () {
            this._appNameChange(this.appNameValue);
        },
        _appNameChange: function (value) {
            this.$appName.addClass('changing');
            setTimeout(function () {
                wlyNavbar.$appName.html(value);
                wlyNavbar.$appName.removeClass('changing');
            }, 100);
        },
        // Bg 设置
        bgSet: function (backgroundValue) {
            this.$container.css('background', backgroundValue);
        },
        // Bg 复原
        bgRec: function () {
            this.$container.css('background', this.bgValue);
        },
        // actionsSceneBar 设置
        actionsSceneBarSet: function (actionsBar, actionsBarClick) {
            this.$actionsBarScene.html(actionsBar);
            this.$actionsBarMain.hide();
            this.$actionsBarScene.show();
            if (actionsBarClick !== undefined) {
                this.$actionsBarScene.find('a').click(actionsBarClick);
            }
        },
        // actionsSceneBar
        actionsSceneBarRec: function () {
            this.$actionsBarScene.hide();
            this.$actionsBarMain.show();
        },
        // NavBar Scene Change
        sceneChange: function (title, bg, actionsBar, actionsBarClick) {
            if (title !== undefined)
                this.appNameSet(title);
            if (bg !== undefined)
                this.bgSet(bg);
            if (actionsBar !== undefined)
                this.actionsSceneBarSet(actionsBar, actionsBarClick);
        },
        // NavBar Scene Rec
        sceneRec: function () {
            this.appNameRec();
            this.bgRec();
            this.actionsSceneBarRec();
        },
        // NavBar 隐藏
        hide: function () {
            this.$container.addClass('navbar-hide');
            this.isShow = false;
        },
        // NavBar 显示
        show: function () {
            this.$container.removeClass('navbar-hide');
            this.isShow = true;
        }
    };

    // Sidebar 操作
    window.wlySidebar = {
        sSidebar: '.sidebar',
        isShow: true,
        onAfterToggle: [],

        init: function () {
            this.$sidebar = $(this.sSidebar);

            $('.sidebar-toggle-btn').click(function () {
                window.wlySidebar.toggle();
            });

            // 移动设备默认隐藏
            if ($(document).width() < 800) {
                this.hide();
            }
        },
        toggle: function () {
            this.isShow ? this.hide() : this.show();

            // 执行已绑定的事件
            if (this.onAfterToggle.length > 0) {
                for (var i in this.onAfterToggle)
                    this.onAfterToggle[i]();
            }
        },
        hide: function () {
            this.$sidebar.addClass('sidebar-hide');
            wlyWrap.changeMainContentAreaFull(true);
            this.isShow = false;
        },
        show: function () {
            this.$sidebar.removeClass('sidebar-hide');
            wlyWrap.changeMainContentAreaFull(false);
            this.isShow = true;
        }
    };
    // WlyQry > wlyTable
    window.wlyTable = {
        $title: '[data-wlytable="title"]',
        $container: '[data-toggle="wlyTable"]',
        $header: '[data-toggle="wlyTable"] .wly-table-header',
        $body: '[data-toggle="wlyTable"] .wly-table-body',
        $bodyItems: '[data-wlytable-item-id]',
        $pagination: '[data-toggle="wlyTable"] .wly-table-pagination',
        $loading: '[data-wlytable="loading"]',
        $fullscreenBtn: '[data-wly-toggle="wlyTableFullScreen"]',
        accessLoad: true,
        pagePer: 50,
        // 初始化
        init: function () {
            var data = wlyTableConfig.data;
            this.pagePer = wlyTableConfig.data.pagination.pagePer;
            // 导入数据
            this.putData(data);
            // 窗口大小发生改变
            $(window).resize(function () {
                wlyTable.screenFit();
            });
            // 设置全屏按钮功能
            $(this.$fullscreenBtn).click(function () {
                wlyTable.screenFull();
            });
            // 设置打印按钮功能
            $('[data-wly-toggle="wlyTablePrint"]').click(function () {
                $(wlyTable.$body).find('table').css('margin-top', '').print({
                    globalStyles: true,
                    mediaPrint: false,
                    iframe: false,
                    prepend: '<h2 style="text-align: center;margin-bottom: 20px">' + $(wlyTable.$title).text() + '</h2>'
                });
                wlyTable.screenFit();
            });
        },
        // 屏幕适配
        screenFit: function () {
            // 设置整个悬浮表格的高度
            // container.height($(window).height()-533);

            var tableHeight = getHeight();
            var isFullScreen = !wlyNavbar.isShow;
            if (isFullScreen) tableHeight = tableHeight + $('.main-navbar').outerHeight(true);
            $(this.$container).css('height', tableHeight + 'px');
            if (!isFullScreen)
                $(this.$container).css('padding-bottom', $(this.$pagination).outerHeight());
            else
                $(this.$container).css('padding-bottom', $(this.$pagination).outerHeight() - $(this.$header).outerHeight() + 20);
            // 设置悬浮样式
            $(this.$body + ' table').css('margin-top', '-' + $(this.$body + ' table thead').outerHeight() - 1);
            // 获取 body table thead tr 中每个 th 对象
            $.each($(this.$body).find('table > thead > tr:first-child:not(.no-records-found) > *'), function (i, item) {
                // 逐个设置 head table 中每个 th 的宽度 === body th 的宽度
                $(wlyTable.$header).find('table > thead th:nth-child(' + parseInt(i + 1) + ')').width($(item).width());
            });
            $(this.$header + ' table').width($(this.$body + ' table').outerWidth(true) - 2); // minus the 2px border-width
            // $(this.$header).css('margin-right', $(this.$body).outerWidth(true)-$(this.$body+' table').outerWidth(true));
        },
        // 全屏
        isScreenFull: false,
        screenFull: function (needScreenFull) {
            var fullScreenClass = 'card--fullscreen';

            if (typeof (needScreenFull) === 'boolean') {
                if (needScreenFull) {
                    if (this.isScreenFull) return;
                    $(this.$container).parents('.card').addClass(fullScreenClass);
                    this.isScreenFull = true;
                } else {
                    if (!this.isScreenFull) return;
                    $(this.$container).parents('.card').removeClass(fullScreenClass);
                    this.isScreenFull = false;
                }
            } else {
                if (this.isScreenFull) {
                    $(this.$container).parents('.card').removeClass(fullScreenClass);
                    this.isScreenFull = false;
                } else {
                    $(this.$container).parents('.card').addClass(fullScreenClass);
                    this.isScreenFull = true;
                }
            }

            // 修改按钮图标
            $(this.$fullscreenBtn).removeClass('zmdi-fullscreen').removeClass('zmdi-fullscreen-exit');
            if (this.isScreenFull) {
                $(this.$fullscreenBtn).html('<i class="zmdi zmdi-fullscreen-exit"></i> <span>退出全屏</span>');
                wlyNavbar.hide();
            } else {
                $(this.$fullscreenBtn).html('<i class="zmdi zmdi-fullscreen"></i> <span>全屏显示</span>');
                wlyNavbar.show();
            }
            // 表格适配屏幕
            this.screenFit();
        },
        // 导入数据
        putData: function (data) {
            wlyTableConfig.data = data;
            var _this = this;

            // 设置 Table Title
            $(this.$title).html(data.dataTitle + ' - ' + data.dataSubtitle + ' <span style="font-size: 13px;vertical-align: bottom;">' + data.dataSubtitleB + '</span>');

            var dataFieldList = data.fieldList,
                dataScore = data.score,
                dataPage = data.pagination,
                dataSortBy = data.sortBy;

            // 过滤规则
            var filterList = []; // ["ranking"];
            var isFiledShow = function (fn) {
                if (filterList.indexOf(fn) > -1)
                    return false;
                if (_this.DisplayController.hiddenFields.indexOf(fn) > -1)
                    return false;
                return true;
            };

            var thead = '',
                tbody = '';
            // $thead
            $.each(dataFieldList, function (i, item) {
                if (!isFiledShow(item.FN))
                    return;

                var itemClass = '';
                var itemTitle = sprintf('%s 降序', item.LB);
                if (dataSortBy[item.FN] && parseInt(dataSortBy[item.FN]) === SORT_ASC) {
                    itemClass = 'select sort-asc';
                    itemTitle = sprintf('%s 降序 [当前为升序]', item.LB);
                } else if (dataSortBy[item.FN] && parseInt(dataSortBy[item.FN]) === SORT_DESC) {
                    itemClass = 'select sort-desc';
                    itemTitle = sprintf('%s 升序 [当前为降序]', item.LB);
                }
                var thValue = '';
                switch (item.FN) {
                    /*case 'overall':
                        thValue = item.LB+'<span class="print-s"> /市排名</span>';
                        break;*/

                    default:
                        thValue = item.LB;
                }
                thead += '<th><span data-wlytable-thead="' + item.FN + '" class="' + itemClass + '" title="依 ' + itemTitle + '">' + thValue + '</i></span></th>';
            });
            thead = sprintf('<thead><tr>%s</tr></thead>', thead);
            // $tbody
            $.each(dataScore, function (scoreId, score) {
                tbody += sprintf('<tr data-wlytable-item-id="%s">', scoreId);
                $.each(dataFieldList, function (fieldIndex, fieldItem) {
                    if (!isFiledShow(fieldItem.FN))
                        return;

                    switch (fieldItem.FN) {
                        /*case 'overall':
                            tbody += sprintf('<td data-wlytable-item-fn="%s">%s <span class="ranking" title="市排名 第%s名"><span class="print-s">/</span>%s</span></td>', fieldItem.FN, score[fieldItem.FN], score.ranking, score.ranking);
                            break;*/

                        case 'name':
                            tbody += (
                                '<td data-wlytable-item-fn="' + fieldItem.FN + '">' +
                                '<span class="table-link" onclick="wlyMine.goToCharts(\'' + score.name + '\', \'' + score.school + '\', \'' + score.clas + '\')">' +
                                score[fieldItem.FN] +
                                '</span>' +
                                '</td>'
                            );
                            break;

                        case 'school':
                            tbody += (
                                '<td data-wlytable-item-fn="' + fieldItem.FN + '">' +
                                '<span class="table-link" onclick="wlySearch.queryBySchool($(this).html())">' +
                                score[fieldItem.FN] +
                                '</span>' +
                                '</td>'
                            );
                            break;

                        case 'clas':
                            tbody += (
                                '<td data-wlytable-item-fn="' + fieldItem.FN + '">' +
                                '<span class="table-link" onclick="wlySearch.queryByClass(wlyTableConfig.data.score[' + scoreId + '].school, $(this).html())">' +
                                score[fieldItem.FN] +
                                '</span>' +
                                '</td>'
                            );
                            break;

                        default:
                            if (typeof(fieldItem.INT) !== "undefined" && fieldItem.INT) {
                                tbody += (
                                    '<td data-wlytable-item-fn="' + fieldItem.FN + '">' +
                                    '<span class="table-link" onclick="wlyMine.goToCharts(\'' + score.name + '\', \'' + score.school + '\', \'' + score.clas + '\', \'' + fieldItem.LB + '\')">' +
                                    score[fieldItem.FN] +
                                    '</span>' +
                                    '</td>'
                                );
                            } else {
                                tbody += '<td data-wlytable-item-fn="' + fieldItem.FN + '">' + score[fieldItem.FN] + '</td>';
                            }
                    }


                });
                tbody += '</tr>';
            });
            tbody = '<tbody>' + tbody + '</tbody>';
            // $pagination
            var pageSelect = '';
            // Show Num
            var pageBtnShowNum = 2;
            // First Page Btn
            var firstPageBtn = '';
            if (dataPage.nowPage > 1 + pageBtnShowNum) {
                firstPageBtn += '<a class="paginate-button" data-wlytable-page="1" title="第一页">1</a>';
            }
            // Left Btn
            for (var pageLeftNum = dataPage.nowPage - pageBtnShowNum; pageLeftNum < dataPage.nowPage; pageLeftNum++) {
                if (pageLeftNum < 1) {
                    continue;
                }
                pageSelect += sprintf('<a class="paginate-button" data-wlytable-page="%s">%s</a>', pageLeftNum, pageLeftNum);
            }
            // Current Page Btn
            pageSelect += sprintf('<a class="paginate-button current" data-wlytable-page="%s">%s</a>', dataPage.nowPage, dataPage.nowPage);
            // Right Btn
            for (var pageRightNum = dataPage.nowPage + 1; pageRightNum < dataPage.nowPage + pageBtnShowNum + 1; pageRightNum++) {
                if (pageRightNum > dataPage.lastPage) {
                    continue;
                }
                pageSelect += sprintf('<a class="paginate-button" data-wlytable-page="%s">%s</a>', pageRightNum, pageRightNum);
            }
            // Last Page Btn
            var lastPageBtn = '';
            if (dataPage.nowPage < dataPage.lastPage - pageBtnShowNum) {
                lastPageBtn += sprintf('<a class="paginate-button" data-wlytable-page="%s" title="最后一页">%s</a>', dataPage.lastPage, dataPage.lastPage);
            }
            var pagination = (
                '<div class="paginate-simple">' +
                sprintf('<a class="paginate-button previous %s" title="上一页" data-wlytable-page="%s"></a>', (dataPage.nowPage === 1 ? 'disabled' : ''), dataPage.nowPage - 1) +
                sprintf('<span>%s</span>', pageSelect) +
                sprintf('<a class="paginate-button next %s" title="下一页" data-wlytable-page="%s"></a>', (dataPage.nowPage === dataPage.lastPage ? 'disabled' : ''), dataPage.nowPage + 1) +
                lastPageBtn +
                '</div>'
            );
            // put in the $container
            $(this.$container).html(
                sprintf('<div class="wly-table-header"><table class="table table-striped table-hover">%s</table></div>', thead) +
                sprintf('<div class="wly-table-body"><table class="table table-striped table-hover">%s %s</table></div>', thead, tbody) +
                sprintf('<div class="wly-table-pagination">%s</div>', pagination) +
                sprintf('<div class="wly-table-loading" style="display: none" data-wlytable="loading"><div class="page-loader__spinner" style="display: none"><svg viewBox="25 25 50 50"><circle cx="50" cy="50" r="20" fill="none" stroke-width="2" stroke-miterlimit="10"></circle></svg></div></div>')
            );
            // Bind Thead Sort Btn
            $('[data-wlytable-thead]').bind('click', function () {
                var fn = $(this).data('wlytable-thead');
                var sortType = SORT_DESC;
                if ($(this).hasClass('sort-desc')) {
                    sortType = SORT_ASC;
                }
                wlyTable.reqData({'sortBy': fn, 'sortType': sortType, 'page': 1}, true);
            });
            // Bind Page Controller Each Btn Click Event
            $('[data-wlytable-page]').bind('click', function () {
                var $this = $(this),
                    pageNum = parseInt($this.data('wlytable-page'));
                if ($(this).hasClass('disabled') || pageNum > dataPage.lastPage || pageNum === dataPage.nowPage || pageNum <= 0) return;
                wlyTable.reqData({'page': pageNum}, true);
            });
            // 适配屏幕
            this.screenFit();
            // $header 跟着 $body 滚动
            $(this.$body).scroll(function (e) {
                $(wlyTable.$header).scrollLeft($(this).scrollLeft());
                $(wlyTable.$header).scrollTop($(this).scrollTop());
                // 滚动时全屏
                wlyTable.screenFull(true);
            });
            // 显示 $container
            $(this.$container).css('opacity', 1);
        },
        // Ajax 请求新数据
        reqData: function (dataObj, extendCurrentParms, hiddenFields) {
            if (!this.accessLoad)
                return;

            var parGet = {};
            if (!!extendCurrentParms) {
                parGet = getRequest();
                delete parGet.timestamp;
            }

            if (!hiddenFields && !extendCurrentParms) {
                this.DisplayController.hiddenFields = [];
            }

            if (!!hiddenFields) {
                // 清空字段隐藏配置 显示全部字段
                this.DisplayController.hiddenFields = [];
                for (var filedKey in hiddenFields) {
                    this.DisplayController.hiddenFields.push(hiddenFields[filedKey])
                }
            }

            var reqObj = $.extend(parGet, {'pagePer': this.pagePer}, dataObj, {'id': wlyTableDataId});

            $.ajax({
                type: 'GET',
                url: '/',
                dataType: "json",
                data: $.extend({'timestamp': new Date().getTime()}, reqObj), // 请求不同的页面，不同的 header 才会生效
                beforeSend: function () {
                    wlyTable.accessLoad = false;
                    wlyTable.loadingSet();
                },
                success: function (json) {
                    if (json.success) {
                        // setTimeout(function() {
                        wlyTable.loadingUnset();
                        // 请求成功
                        wlyTable.putData(json.data);
                        // console.log(req);
                        setUrlByParObj(reqObj);
                        // },3000);
                    } else {
                        alert(json.msg);
                        wlyTable.loadingUnset();
                    }
                    wlyTable.accessLoad = true;
                },
                error: function () {
                    // 请求失败
                    wlyTable.accessLoad = true;
                    wlyTable.loadingUnset();
                    alert('请求失败，网络异常');
                    throw ('wlyTable: 请求失败，网络异常');
                }
            });
        },
        // 加载指示器设置
        loadingSet: function () {
            $(wlyTable.$loading).fadeIn();
            setTimeout(function () {
                if ($(wlyTable.$loading).css('display') === 'none')
                    return;

                $(wlyTable.$loading).find('.page-loader__spinner').fadeIn();
            }, 400);
        },
        // 加载指示器解除
        loadingUnset: function () {
            $(wlyTable.$loading).find('.page-loader__spinner').fadeOut();
            $(wlyTable.$loading).fadeOut();
            $(wlyTable.$body).css('filter', '');
        },
        // 显示数据统计面板
        showDataCounter: function () {
            var dialog = new window.wlyDialogBuilder('wly-table-data-counter', '数据统计');
            var dialogBodyElem = dialog.setBody(
                '<span class="dialog-label">数据 "' + wlyTableConfig.data.dataSubtitle + '" 平均值</span>'
            );

            var avgs = wlyTableConfig.data.dataAvg;
            for (var key in avgs) {
                var avgItem = avgs[key];
                var itemDom = $(
                    '<span class="data-item">' +
                    '<span class="data-name">' + avgItem[1].LB + '</span>' +
                    '<span class="data-value">' + avgItem[0] + '</span>' +
                    '</span>'
                ).appendTo(dialogBodyElem);
            }
        },

        DisplayController: {
            hiddenFields: [],
            show: function () {
                var dialog = new window.wlyDialogBuilder('display-controller', '表格显示调整');
                var dialogBodyElem = dialog.setBody(
                    '<span class="dialog-label">点按下列方块来 显示 / 隐藏 字段</span>' +
                    '<div class="field-list"></div>' +
                    '<span class="dialog-label">每页显示项目数量 （数字不宜过大）</span>' +
                    '<div class="page-per-show">' +
                    '<input type="number" class="page-per-show-input" placeholder="每页显示数" min="1" value="' + wlyTable.pagePer + '" />' +
                    '</div>' +

                    '<span class="dialog-label">表格字体大小调整</span>' +
                    '<div class="table-font-size-control">' +
                    '<span class="font-size-minus">-</span>' +
                    '<span class="font-size-value"></span>' +
                    '<span class="font-size-plus">+</span>' +
                    '</div>'
                );

                var fields = window.wlyTableConfig.data.fieldList;
                var _this = this;
                for (var key in fields) {
                    var field = fields[key];
                    var fieldItemDom = $('<span class="field-item" data-field-name="' + field.FN + '">' + field.LB + '</span>');
                    if ($.inArray(field.FN, this.hiddenFields) < 0) {
                        fieldItemDom.addClass('active');
                    }
                    fieldItemDom.click(function () {
                        var fieldName = $(this).attr('data-field-name');
                        if ($(this).hasClass('active')) {
                            // 字段隐藏
                            _this.fieldHide(fieldName);
                            $(this).removeClass('active');
                        } else {
                            // 字段显示
                            _this.fieldShow(fieldName);
                            $(this).addClass('active');
                        }
                    });
                    fieldItemDom.appendTo(dialogBodyElem.find('.field-list'));
                }
                var inputNumBackup = dialogBodyElem.find('.page-per-show-input').val();
                dialogBodyElem.find('.page-per-show-input').change(function () {
                    var value = $(this).val();
                    if (isNaN(Number(value))) {
                        alert('只能输入数字');
                        return;
                    }
                    wlyTable.pagePer = value;
                    wlyTable.reqData({'page': '1'}, true);
                });
                var currentTableFontSize = Number($('.table').first().css('font-size').replace(/px$/, ''));
                var changeTableFontSize = function (value) {
                    var tableFontSizeStyleDom = $('#TableFontSize');
                    if (tableFontSizeStyleDom.length === 0) {
                        tableFontSizeStyleDom = $('<style id="TableFontSize"></style>').appendTo('head');
                    }
                    tableFontSizeStyleDom.html('.table {font-size: ' + value + 'px;}');
                    currentTableFontSize = value;
                    dialogBodyElem.find('.font-size-value').text(currentTableFontSize);
                };
                dialogBodyElem.find('.font-size-value').text(currentTableFontSize);
                dialogBodyElem.find('.font-size-minus').click(function () {
                    changeTableFontSize(currentTableFontSize - 2);
                    wlyTable.screenFit();
                });
                dialogBodyElem.find('.font-size-plus').click(function () {
                    changeTableFontSize(currentTableFontSize + 2);
                    wlyTable.screenFit();
                });
            },
            // 字段隐藏
            fieldHide: function (fieldName) {
                if ($.inArray(fieldName, this.hiddenFields) > -1) return;
                this.hiddenFields.push(fieldName);
                wlyTable.putData(wlyTableConfig.data);
            },
            // 字段显示
            fieldShow: function (fieldName) {
                if ($.inArray(fieldName, this.hiddenFields) < 0) return;
                arrayRemoveByValue(this.hiddenFields, fieldName);
                wlyTable.putData(wlyTableConfig.data);
            }
        }
    };
    // WlyQry > wlySearch
    window.wlySearch = {
        $container: 'body .wly-search-panel',
        // 显示面板
        showPanel: function () {
            // 若已存在
            if ($(this.$container).length !== 0) return false;

            // 内容
            var container = $(
                '<div class="wly-search-panel anim-fade-in">' +
                '<div class="search-panel-inner"></div>' +
                '</div>'
            ).appendTo('body');

            var searchPanelInnerDom = container.find('.search-panel-inner');

            // 通过姓名查询成绩
            var searchByNameDom = $(
                '<div class="card container">' +
                '<div class="card-block">' +
                '<form class="search-form" data-wlysearch="query-by-name-form">' +
                '<input class="search-input" type="text" name="queryData" placeholder="在此输入学生姓名..." required="required" autocomplete="off" />' +
                '<button type="submit" class="search-btn"><i class="zmdi zmdi-search"></i></button>' +
                '</form>' +
                '</div>' +
                '</div>'
            ).appendTo(searchPanelInnerDom);

            // 通过类别来查成绩
            var searchByCategory = $(
                '<div class="card container">' +
                '<div class="card-block">' +
                '<form class="search-form query-by-category-form" data-wlysearch="query-by-category-form">' +
                '<select name="queryData[school]"><option value="">选择学校</option></select>' +
                '<span class="middle-text">AND</span>' +
                '<select name="queryData[clas]"><option value="">所有班级</option></select>' +
                '<button type="submit" class="search-btn"><i class="zmdi zmdi-search"></i></button>' +
                '</form>' +
                '</div>' +
                '</div>'
            );
            searchByCategory.hide() // 默认是隐藏起来的
                .appendTo(searchPanelInnerDom);

            // 之后的工具条
            var toolbarDom = $(
                '<div class="search-toolbar">' +
                '<button type="button" class="action-btn" data-wlysearch="toggle-query-mode">根据 学校班级 查询</button>' +
                '<button type="button" class="action-btn" data-wlysearch="view-all-data">查看全市成绩</button>' +
                '<span class="copyright">' + wlyTableConfig.sign + ' <a href="https://github.com/ZNEIAT" target="_blank">ZNEIAT/QWQUERY</a></span>' +
                '</div>'
            ).appendTo(searchPanelInnerDom);

            /** Begin 样式代码 **/
            // 显示面板的动画
            setTimeout(function () {
                $(container).find('.search-panel-inner').addClass('show-panel');
                setTimeout(function () {
                    $(container).css('overflow', 'auto');
                }, 400);
            }, 400);
            // 导航栏
            wlyNavbar.sceneChange('搜索', '#03A9F4', '<li><a data-scene-action-btn="closeSearch"><i class="zmdi zmdi-close"></i></a></li>', function () {
                var data = $(this).data('scene-action-btn');
                switch (data) {
                    case 'closeSearch':
                        wlySearch.removePanel();
                        break;
                }
            });
            // Body
            $('body').css('overflow', 'hidden');
            /** End 样式代码 **/

            // 通过姓名查询成绩
            searchByNameDom.find('.search-input').focus();
            searchByNameDom.find('[data-wlysearch="query-by-name-form"]').submit(function () {
                var input = $(this).find('[name="queryData"]');
                var value = $.trim(input.val());
                if (value === '') {
                    input.focus();
                    return false;
                }
                wlySearch.removePanel();
                wlyTable.reqData({'queryType': 'name', 'queryData': value, 'page': '1'});
                return false;
            });

            // 通过类别来查成绩
            // Bind The Form Submit Event -> just for check the input Value...
            var categorySchoolInputDom = $('[name="queryData[school]"]');
            var categoryClasInputDom = $('[name="queryData[clas]"]');
            // categorySchoolInputDom 数据供应
            $.ajax({
                type: 'GET',
                url: '/site/get-all-category',
                dataType: "json",
                data: {'id': wlyTableDataId},
                beforeSend: function () {
                    categorySchoolInputDom.html('<option value="" selected>数据加载中...</option>');
                },
                success: function (json) {
                    if (json.success) {
                        var data = json.data;
                        var htmlCode = '<option value="" selected>选择学校</option>';
                        for (var o in data) {
                            htmlCode += '<option value="' + data[o] + '">' + data[o] + '</option>';
                        }
                        categorySchoolInputDom.html(htmlCode);
                    } else {
                        categorySchoolInputDom.html('<option value="" selected>数据加载失败</option>');
                    }
                },
                error: function () {
                    categorySchoolInputDom.html('<option value="" selected>数据加载失败</option>');
                    throw ('/site/get-all-category: 请求失败，网络异常');
                }
            });
            // categoryClasInputDom 数据供应
            categorySchoolInputDom.change(function () {
                var schoolName = $(this).val();
                $.ajax({
                    type: 'GET',
                    url: '/site/get-all-category',
                    dataType: "json",
                    data: {'id': wlyTableDataId, 'school': schoolName},
                    beforeSend: function () {
                        categoryClasInputDom.html('<option value="" selected>数据加载中...</option>');
                    },
                    success: function (json) {
                        if (json.success) {
                            var data = json.data;
                            var htmlCode = '<option value="" selected>所有 ' + schoolName + ' 的班级</option>';
                            for (var o in data) {
                                htmlCode += '<option value="' + data[o] + '">' + data[o] + '</option>';
                            }
                            categoryClasInputDom.html(htmlCode);
                        } else {
                            categoryClasInputDom.html('<option value="" selected>数据加载失败</option>');
                        }
                    },
                    error: function () {
                        categoryClasInputDom.html('<option value="" selected>数据加载失败</option>');
                        throw ('/site/get-all-category|categoryClasInputDom: 请求失败，网络异常');
                    }
                });
            });
            searchByCategory.find('[data-wlysearch="query-by-category-form"]').submit(function () {
                if (categorySchoolInputDom.val().length === 0) {
                    categorySchoolInputDom.focus();
                    return false;
                }

                wlySearch.removePanel();
                if (categoryClasInputDom.val().length === 0) {
                    wlySearch.queryBySchool(categorySchoolInputDom.val());
                } else {
                    wlySearch.queryByClass(categorySchoolInputDom.val(), categoryClasInputDom.val());
                }

                return false;
            });

            // 工具条功能：分类查询模式切换
            toolbarDom.find('[data-wlysearch="toggle-query-mode"]').click(function () {
                if (searchByCategory.css('display') === 'none') {
                    searchByCategory.show();
                    searchByNameDom.hide();
                    $(this).text('根据 姓名 查询');
                } else {
                    searchByCategory.hide();
                    searchByNameDom.show();
                    $(this).text('根据 学校班级 查询');
                }
            });
            // 工具条功能：全市成绩
            toolbarDom.find('[data-wlysearch="view-all-data"]').click(function () {
                wlySearch.removePanel();
                wlyTable.reqData({});
            });

            // 导航栏显示隐藏
            if (!wlyNavbar.isShow) {
                this.removePanelAndNavHide = true;
                wlyNavbar.show();
            } else {
                this.removePanelAndNavHide = false;
            }
        },
        removePanelAndNavHide: false,
        // 移除面板
        removePanel: function () {
            // 若不存在
            if ($(this.$container).length === 0)
                return false;
            // 内容
            $(this.$container).remove();
            // Body
            $('body').css('overflow', '');
            // 导航栏
            wlyNavbar.sceneRec();
            if (this.removePanelAndNavHide) {
                wlyNavbar.hide();
            }
        },
        queryBySchool: function (schoolName) {
            wlyTable.reqData({
                'queryType': 'school',
                'queryData_school': schoolName,
                'page': '1'
            }, false, ['school']);
        },
        queryByClass: function (schoolName, className) {
            wlyTable.reqData({
                'queryType': 'class',
                'queryData_school': schoolName,
                'queryData_class': className,
                'page': '1'
            }, false, ['school', 'clas']);
        }
    };

    // 对话框创建器
    window.wlyDialogBuilder = function (_className, _title) {
        var dialogElem = $(
            '<div class="wly-table-action-dialog anim-fade-in ' + _className + '">' +
            '<div class="wly-table-action-dialog-inner">' +

            '<div class="dialog-title">' +
            '<span class="title-text"></span>' +
            '<span data-dialog-func="close" class="close-btn zmdi zmdi-close"></span>' +
            '</div>' +

            '<div class="dialog-body"></div>' +

            '</div>' +
            '</div>'
        ).appendTo('body');

        this.getElem = function () {
            return dialogElem;
        };
        this.setTile = function (title) {
            dialogElem.find('.title-text').html(title);
        };
        this.setBody = function (bodyContent) {
            var bodyElem = this.getBodyElem();
            return bodyElem.html(bodyContent);
        };
        this.getBodyElem = function () {
            return dialogElem.find('.dialog-body');
        };

        this.setTile(_title);
        dialogElem.find('[data-dialog-func="close"]').click(function () {
            dialogElem.remove();
        });

    };

    // 保存数据为电子表格
    window.wlyTableDataSave = {
        show: function () {
            var nowPage = wlyTableConfig.data.pagination.nowPage;
            var lastPage = wlyTableConfig.data.pagination.lastPage;

            var dialog = new window.wlyDialogBuilder('wly-table-data-save', '保存数据为电子表格');
            var dialogBodyElem = dialog.setBody(
                '<span class="dialog-label">保存 ' + wlyTableConfig.data.dataTitle + ' 的 ' + wlyTableConfig.data.dataSubtitle + ' 为电子表格</span>' +
                '<span class="dialog-btn" data-dialog-func="save-now">保存 仅第 ' + nowPage + ' 页 数据</span>' +
                '<span class="dialog-btn" data-dialog-func="save-now-noPaging">保存 第 ' + nowPage + '~' + lastPage + ' 页 数据</span>' +
                '<span class="dialog-label">保存 ' + wlyTableConfig.data.dataTitle + ' 全部数据为电子表格</span>' +
                '<span class="dialog-btn" data-dialog-func="save-noPaging">保存 全市成绩</span>'
            );

            dialogBodyElem.find('[data-dialog-func="save-now"]').click(function () {
                location.href = '/?' + jQuery.param($.extend(getRequest(), {"saveFile": 'y'}));
            });
            dialogBodyElem.find('[data-dialog-func="save-now-noPaging"]').click(function () {
                location.href = '/?' + jQuery.param($.extend(getRequest(), {"saveFile": 'y', 'saveMode': 'noPaging'}));
            });
            dialogBodyElem.find('[data-dialog-func="save-noPaging"]').click(function () {
                location.href = '/?' + jQuery.param({'id': wlyTableDataId, 'saveFile': 'y', 'saveMode': 'noPaging'});
            });
        }
    };

    // 个人信息显示
    window.wlyMine = {
        goToCharts: function (name, schoolName, className, filedName) {
            var params = {};
            params.queryName = name;
            params.querySchool = schoolName;
            params.queryClass = className;
            if (typeof(filedName) !== "undefined")
                params.onlyFiledName= filedName;
            window.location.href = '/site/charts' + parseParamUrl(params);
        }
    };

    console.log("%c qwquery %c 作者 ZNEIAT<1149527164@qq.com> %c 项目地址：https://github.com/Zneiat/qwquery ","color: #FFF; background: #1DAAFF; padding:5px 0;","color: #FFF; background: #656565; padding:5px 0;", "color: #656565; background: #FFF; padding:5px 0;");
    console.log('%c未经允许代码和衍生品不得用于商业用途，侵权必究', 'color: red;padding:10px 0;');

    /**
     * 删除数组中的指定值
     */
    function arrayRemoveByValue(arr, val) {
        for (var i = 0; i < arr.length; i++) {
            if (arr[i] === val) {
                arr.splice(i, 1);
                break;
            }
        }
    }

    /**
     * 字符串格式化
     */
    function sprintf(str) {
        var args = arguments,
            flag = true,
            i = 1;

        str = str.replace(/%s/g, function () {
            var arg = args[i++];

            if (typeof arg === 'undefined') {
                flag = false;
                return '';
            }
            return arg;
        });
        return flag ? str : '';
    }

    /**
     * 获取 URL 中的请求参数
     */
    function getRequest() {
        var strs,
            url = location.search, //获取url中"?"符后的字串
            theRequest = {};
        if (url.indexOf("?") !== -1) {
            var str = url.substr(1);
            strs = str.split("&");
            for (var i = 0; i < strs.length; i++) {
                theRequest[strs[i].split("=")[0]] = decodeURIComponent(strs[i].split("=")[1]);
            }
        }
        return theRequest;
    }

    /**
     * 对象转 URL 参数
     */
    function parseParamUrl(param) {
        var string = '?';
        $.each(param, function (key, val) {
            string += key + "=" + encodeURIComponent(val) + "&";
        });
        string = string.substring(0, string.length - 1);
        return string;
    }

    /**
     * 用参数对象设置当前页的 URL
     */
    function setUrlByParObj(paramObj) {
        var parStr = parseParamUrl(paramObj);
        // console.log(parStr);
        History.pushState(null, document.title, parStr);
    }

    /**
     * 合并两个对象
     */
    function objExtend(obj1, obj2) {
        for (var key in obj2) {
            if (obj1.hasOwnProperty(key)) continue; // 有相同的属性则略过
            obj1[key] = obj2[key];
        }
        return obj1;
    }
})();
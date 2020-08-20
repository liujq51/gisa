<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <title>AdminLTE 2 | with iframe</title>
    <!-- Tell the browser to be responsive to screen width -->
    <meta content="width=device-width, initial-scale=1, maximum-scale=1, user-scalable=no" name="viewport">

    <!-- Bootstrap 3.3.6 -->
    <link rel="stylesheet" href="/static/adminlte/bootstrap/css/bootstrap.min.css">
    <!-- Font Awesome -->
    <link rel="stylesheet" href="/static/adminlte/dist/css/font-awesome.min.css">
    <!-- Ionicons -->
    <link rel="stylesheet" href="/static/adminlte/dist/css/ionicons.min.css">
    <!-- Theme style -->
    <link rel="stylesheet" href="/static/adminlte/dist/css/AdminLTE.min.css">
    <link rel="stylesheet" href="/static/adminlte/dist/css/skins/all-skins.min.css">

    <style type="text/css">
        html {
            overflow: hidden;
        }
    </style>
    <!-- HTML5 Shim and Respond.js IE8 support of HTML5 elements and media queries -->
    <!-- WARNING: Respond.js doesn't work if you view the page via file:// -->
    <!--[if lt IE 9]>
    <script src="/static/adminlte/plugins/ie9/html5shiv.min.js"></script>
    <script src="/static/adminlte/plugins/ie9/respond.min.js"></script>
    <![endif]-->
</head>
<body class="hold-transition skin-black sidebar-mini fixed">
<div class="wrapper">

    {{.HtmlHead}}

    {{.HtmlSidebar}}

    <!-- Content Wrapper. Contains page content -->
    <div class="content-wrapper" id="content-wrapper" style="min-height: 421px;">
        <!--bootstrap tab风格 多标签页-->
        <div class="content-tabs">
            <button class="roll-nav roll-left tabLeft" onclick="scrollTabLeft()">
                <i class="fa fa-backward"></i>
            </button>
            <nav class="page-tabs menuTabs tab-ui-menu" id="tab-menu">
                <div class="page-tabs-content" style="margin-left: 0px;">
                </div>
            </nav>
            <button class="roll-nav roll-right tabRight" onclick="scrollTabRight()">
                <i class="fa fa-forward" style="margin-left: 3px;"></i>
            </button>
            <div class="btn-group roll-nav roll-right">
                <button class="dropdown tabClose" data-toggle="dropdown">
                页签操作<i class="fa fa-caret-down" style="padding-left: 3px;"></i>
                </button>
                <ul class="dropdown-menu dropdown-menu-right" style="min-width: 128px;">
                    <li><a class="tabReload" href="javascript:refreshTab();">刷新当前</a></li>
                    <li><a class="tabCloseCurrent" href="javascript:closeCurrentTab();">关闭当前</a></li>
                    <li><a class="tabCloseAll" href="javascript:closeOtherTabs(true);">全部关闭</a></li>
                    <li><a class="tabCloseOther" href="javascript:closeOtherTabs();">除此之外全部关闭</a></li>
                </ul>
            </div>
            <button class="roll-nav roll-right fullscreen" onclick="App.handleFullScreen()"><i class="fa fa-arrows-alt"></i></button>
        </div>
        <div class="content-iframe " style="background-color: #ffffff; ">
            <div class="tab-content " id="tab-content">
            </div>
        </div>
    </div>
    <!-- /.content-wrapper -->
   {{.HtmlFooter}}
    <div class="control-sidebar-bg"></div>
</div>
<!-- ./wrapper -->

    <!-- jQuery 2.2.3 -->
    <script src="/static/adminlte/plugins/jQuery/jquery-2.2.3.min.js"></script>
    <!-- Bootstrap 3.3.6 -->
    <script src="/static/adminlte/bootstrap/js/bootstrap.min.js"></script>
    <!-- Slimscroll -->
    <script src="/static/adminlte/plugins/slimScroll/jquery.slimscroll.min.js"></script>
    <!-- FastClick -->
    <script src="/static/adminlte/plugins/fastclick/fastclick.js"></script>
    <!-- AdminLTE App -->
    <script src="/static/adminlte/dist/js/app.js"></script>
    <!-- AdminLTE for demo purposes -->
    <script src="/static/adminlte/dist/js/demo.js"></script>

    <!--tabs-->
    <script src="/static/adminlte/dist/js/app_iframe.js"></script>

    <!--<script src="/static/adminlte/dist/js/jquery.blockui.min.js"></script>
    <script src="/static/adminlte/dist/js/appx.js"></script>
    <script src="/static/adminlte/dist/js/bootstrap-tab.js"></script>
    <script src="/static/adminlte/dist/js/sidebarMenu.js"></script>-->

    <script type="text/javascript">
        $(function () {
    //        console.log(window.location);
            App.setbasePath("./");
            App.setGlobalImgPath("static/adminlte/dist/img/");

            addTabs({
                id: '99999',
                title: '欢迎页',
                close: false,
                url: 'dashboard',
                urlType: "relative"
            });

            App.fixIframeCotent();
            eval("var menus="+{{.Menus}});
            $('.sidebar-menu').sidebarMenu({data: menus});

            // 动态创建菜单后，可以重新计算 SlimScroll
            // $.AdminLTE.layout.fixSidebar();
        });
    </script>
</body>
</html>
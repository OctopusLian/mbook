$(function () {
    window.addDocumentModalFormHtml = $(this).find("form").html();
    window.editor = editormd("docEditor", {
        atLink    : false,    // disable @link
        emailLink : false,    // disable email address auto link
        emoji : false,       // Support Github emoji, Twitter Emoji(Twemoji), fontAwesome, Editor.md logo emojis.
        width : "100%",
        height : "100%",
        path : "/static/editor.md/lib/",
        toolbar : true,
        placeholder: "本编辑器支持Markdown编辑，左边编写，右边预览",
        imageUpload: true,
        imageFormats: ["jpg", "jpeg", "gif", "png", "JPG", "JPEG", "GIF", "PNG"],
        imageUploadURL: window.imageUploadURL ,
        toolbarModes : "full",
        fileUpload: true,
        fileUploadURL : window.fileUploadURL,
        taskList : true,
        tex  : true,
        flowChart : true,
        sequenceDiagram: true,
        htmlDecode : "style,script,iframe,title,onmouseover,onmouseout,style",
        lineNumbers : true,
        tocStartLevel : 1,
        tocm : true,
        // watch: true,
        saveHTMLToTextarea : true,
        onload : function() {
            this.hideToolbar();
            var keyMap = {
                "Ctrl-S": function(cm) {
                    saveDocument(false);
                },
                "Cmd-S" : function(cm){
                    saveDocument(false);
                },
                "Ctrl-A": function(cm) {
                    cm.execCommand("selectAll");
                }
            };
            this.addKeyMap(keyMap);

            var $select_node_id = window.treeCatalog.get_selected();


            if($select_node_id) {
                var $select_node = window.treeCatalog.get_node($select_node_id[0])
                if ($select_node) {
                    $select_node.node = {
                        id: $select_node.id
                    };

                    loadDocument($select_node);
                }
            }
            uploadImage("docEditor",function ($state, $res) {
                if($state === "before"){
                    return layer.load(1, {
                        shade: [0.1,'#fff'] //0.1透明度的白色背景
                    });
                }else if($state === "success"){
                    if($res.errcode === 0) {
                        var value = '![](' + $res.url + ')';
                        window.editor.insertValue(value);
                    }
                }
            });
        },
        onchange : function () {
            resetEditorChanged(true);
        }
    });





    // window.editor.tocDropdown=true;

    /**
     * 实现标题栏操作
     */
    $("#editormd-tools").on("click","a[class!='disabled']",function () {
       var name = $(this).find("i").attr("name");
       if(name === "attachment"){
           $("#uploadAttachModal").modal("show");
       }else if(name === "history"){
           window.documentHistory();
       }else if(name === "save"){
            saveDocument(false);
       }else if(name === "template"){
           $("#documentTemplateModal").modal("show");
       }else if(name === "sidebar"){
            $("#manualCategory").toggle(0,"swing",function () {

                var $then = $("#manualEditorContainer");
                var left = parseInt($then.css("left"));
                if(left > 0){
                    window.editorContainerLeft = left;
                    $then.css("left","0");
                }else{
                    $then.css("left",window.editorContainerLeft + "px");
                }
                window.editor.resize();
            });
       }else if(name === "release"){
            if(Object.prototype.toString.call(window.documentCategory) === '[object Array]' && window.documentCategory.length > 0){
                if($("#markdown-save").hasClass('change')) {
                    var comfirm_result = confirm("编辑内容未保存，需要保存吗？")
                    if(comfirm_result) {
                        saveDocument();
                    }
                }
                $.ajax({
                    url : window.releaseURL,
                    data :{"identify" : window.book.identify },
                    type : "post",
                    dataType : "json",
                    success : function (res) {
                        if(res.errcode === 0){
                            layer.msg("发布任务已推送到任务队列，稍后将在后台执行。");
                        }else{
                            layer.msg(res.message);
                        }
                    }
                });
            }else{
                layer.msg("没有需要发布的文档")
            }
       }else if(name === "generate"){
           $.ajax({
               url : window.generateURL,
               data :{"identify" : window.book.identify },
               type : "post",
               dataType : "json",
               success : function (res) {
                   layer.msg(res.message);
               }
           }).fail(function () {
               layer.msg("请求失败");
           });
       }else if(name === "tasks") {
           //插入GFM任务列表
           var cm = window.editor.cm;
           var selection = cm.getSelection();

           if (selection === "") {
               cm.replaceSelection("- [x] " + selection);
           }
           else {
               var selectionText = selection.split("\n");

               for (var i = 0, len = selectionText.length; i < len; i++) {
                   selectionText[i] = (selectionText[i] === "") ? "" : "- [x] " + selectionText[i];
               }
               cm.replaceSelection(selectionText.join("\n"));
           }
       }else if(name=="summary"){//目录排序
           var cm = window.editor.cm;
           var selection = cm.getSelection();
           if (selection === "") {
               cm.replaceSelection("<mbook-summary></mbook-summary>" + selection);
           }else {
               var selectionText = selection.split("\n");
               for (var i = 0, len = selectionText.length; i < len; i++) {
                   selectionText[i] = (selectionText[i] === "") ? "" : "<mbook-summary></mbook-summary>" + selectionText[i];
               }
               cm.replaceSelection(selectionText.join("\n"));
           }
       }else if(name=="auto"){//自动生成内容
           var cm = window.editor.cm;
           var selection = cm.getSelection();
           if (selection === "") {
               cm.replaceSelection("<mbook-auto></mbook-auto>" + selection);
           }else {
               var selectionText = selection.split("\n");
               for (var i = 0, len = selectionText.length; i < len; i++) {
                   selectionText[i] = (selectionText[i] === "") ? "" : "<mbook-auto></mbook-auto>" + selectionText[i];
               }
               cm.replaceSelection(selectionText.join("\n"));
           }
       }else if(name=="commit"){//自动生成内容
           var cm = window.editor.cm;
           var selection = cm.getSelection();
           var cursor    = cm.getCursor();
           cm.replaceSelection("<mbook-git>" + selection + "</mbook-git>");
           console.log(cursor.line,cm.lineCount());
           // cm.setCursor(cursor.line, cursor.ch + '<mbook-git>'.length);
       }else if(name=="multi"){//批量创建文档
            $("#ModalMulti").modal("show");
       }else if(name=="spider"){//爬虫采集
            $("#ModalSpider").modal("show");
       }else if(name=="replace"){//全局内容替换
           $("#ModalReplace").modal("show");
       }else{
           var action = window.editor.toolbarHandlers[name];
           if (action !== "undefined") {
               $.proxy(action, window.editor)();
               window.editor.focus();
           }
       }
   }) ;

    /***
     * 加载指定的文档到编辑器中
     * @param $node
     */
    window.loadDocument = function($node) {
        var index = layer.load(1, {
            shade: [0.1,'#fff'] //0.1透明度的白色背景
        });

        $.get(window.editURL + $node.node.id ).done(function (res) {
            layer.close(index);
            resetEditor();
            if(res.errcode === 0){
                window.isLoad = true;
                window.editor.clear();
                window.editor.insertValue(res.data.markdown);
                window.editor.setCursor({line:0, ch:0});
                var node = { "id" : res.data.doc_id,'parent' : res.data.parent_id === 0 ? '#' : res.data.parent_id ,"text" : res.data.doc_name,"identify" : res.data.identify,"version" : res.data.version};
                pushDocumentCategory(node);
                window.selectNode = node;
                pushVueLists(res.data.attach);
                // 设置 history
                if(!window.onpop) history.pushState(node,res.data.doc_name,window.editURI + res.data.identify);
                window.onpop=false;
            }else{
                layer.msg("文档加载失败");
            }
        }).fail(function () {
            layer.close(index);
            layer.msg("文档加载失败");
        });
    };

    /**
     * 保存文档到服务器
     * @param $is_cover 是否强制覆盖
     */
    function saveDocument($is_cover,callback) {
        var index = null;
        var node = window.selectNode;
        var content = window.editor.getMarkdown();
        var html = window.editor.getPreviewedHTML();
        // var html = window.editor.getHTML();
        var version = "";
        var cm = window.editor.cm;

        if(!node){
            layer.msg("获取当前文档信息失败");
            return;
        }
        var doc_id = parseInt(node.id);

        for(var i in window.documentCategory){
            var item = window.documentCategory[i];

            if(item.id === doc_id){
                version = item.version;
                break;
            }
        }
        $.ajax({
            beforeSend  : function () {
                index = layer.load(1, {shade: [0.1,'#fff'] });
            },
            url :  window.editURL,
            data : {"identify" : window.book.identify,"doc_id" : doc_id,"markdown" : content,"html" : html,"cover" : $is_cover ? "yes":"no","version": version},
            type :"post",
            dataType :"json",
            success : function (res) {
                layer.close(index);
                if(res.errcode === 0){
                        resetEditorChanged(false);
                        for(var i in window.documentCategory){
                            var item = window.documentCategory[i];

                            if(item.id === doc_id && res.data!=undefined && res.data.version!=undefined){
                                window.documentCategory[i].version = res.data.version;
                                break;
                            }
                        }
                        if(typeof callback === "function"){
                            callback();
                        }

                        if(res.message=="true"){//刷新页面显示最新的排序
                            location.href=location.origin+location.pathname+"?"+new Date();//刷新当前页面以获取新的排序
                        }else if(res.message=="auto"){
                            $(".jstree-wholerow-clicked").trigger("click");
                        }
                }else if(res.errcode === 6005){
                    var confirmIndex = layer.confirm('文档已被其他人修改确定覆盖已存在的文档吗？', {
                        btn: ['确定','取消'] //按钮
                    }, function(){
                        layer.close(confirmIndex);
                        saveDocument(true,callback);
                    });
                }else{
                    layer.msg(res.message);
                }
            }
        });
    }

    function resetEditor($node) {

    }

    /**
     * 设置编辑器变更状态
     * @param $is_change
     */
    function resetEditorChanged($is_change) {
        if($is_change && !window.isLoad ){
            $("#markdown-save").removeClass('disabled').addClass('change');
        }else{
            $("#markdown-save").removeClass('change').addClass('disabled');
        }
        window.isLoad = false;
    }
    /**
     * 添加顶级文档
     */
    $("#addDocumentForm").ajaxForm({
        beforeSubmit : function () {
            var doc_name = $.trim($("#documentName").val());
            if (doc_name === ""){
                return showError("目录名称不能为空","#add-error-message")
            }
            $("#btnSaveDocument").button("loading");
            return true;
        },
        success : function (res) {
            if(res.errcode === 0){
                var identify=res.data.identify||res.data.doc_id;
                var data = { "id" : res.data.doc_id,'parent' : res.data.parent_id === 0 ? '#' : res.data.parent_id ,"text" : res.data.doc_name+'<small class="text-danger">('+identify+')</small>',"identify" : res.data.identify,"version" : res.data.version};

                var node = window.treeCatalog.get_node(data.id);
                if(node){
                    window.treeCatalog.rename_node({"id":data.id},data.text);

                }else {
                    window.treeCatalog.create_node(data.parent, data);
                    window.treeCatalog.deselect_all();
                    window.treeCatalog.select_node(data);
                }
                pushDocumentCategory(data);
                $("#markdown-save").removeClass('change').addClass('disabled');
                $("#addDocumentModal").modal('hide');
            }else{
                showError(res.message,"#add-error-message")
            }
            $("#btnSaveDocument").button("reset");
        }
    });

    //提交采集
    $("#btnCrawl").click(function (e) {
        e.preventDefault();
        if($(this).hasClass("disabled")) return false;
        var form=$("#ModalSpider form"),action=form.attr("action"),_url=form.find("[name=url]").val();
        if(_url==""){
            form.find("[name=url]").focus();
            layer.msg("内容链接地址不能为空");
        }else{
            $.post(action,form.serialize(),function (res) {
                $("#btnCrawl").addClass("disabled");
                layer.msg("提交成功");
                if (res.errcode==0){
                    layer.msg("采集成功");
                    var cm = window.editor.cm;
                    var selection = cm.getSelection();
                    if (selection === "") {
                        cm.replaceSelection(res.data + selection);
                    }else {
                        var selectionText = selection.split("\n");
                        for (var i = 0, len = selectionText.length; i < len; i++) {
                            selectionText[i] = (selectionText[i] === "") ? "" : res.data + selectionText[i];
                        }
                        cm.replaceSelection(selectionText.join("\n"));
                    }
                    form.find("[name=url]").val("");
                    form.find("[type=reset]").trigger("click");
                }else{
                    layer.msg(res.message);
                }
                $("#btnCrawl").removeClass("disabled");
            });
        }
    });


    //提交替换
    $("#btnReplace").click(function (e) {
        e.preventDefault();
        if($(this).hasClass("disabled")) return false;
        var form=$("#ModalReplace form"),action=form.attr("action"),src=form.find("[name=src]").val();
        if(src==""){
            form.find("[name=src]").focus();
            layer.msg("源内容字符串不能为空");
        }else{
            $("#btnReplace").addClass("disabled");
            $.post(action,form.serialize(),function (res) {
                layer.msg(res.message);
                if (res.errcode==0){
                    setTimeout(function () {
                        location.reload();
                    },1000)
                }
            });
            $("#btnReplace").removeClass("disabled");
        }
    });

    $("#btnMulti").click(function (e) {
        e.preventDefault();
        if($(this).hasClass("disabled")) return false;
        var form=$("#ModalMulti form"),action=form.attr("action"),content=form.find("[name=content]").val();
        if(content==""){
            form.find("[name=content]").focus();
            layer.msg("请添加章节");
        }else{
            $.post(action,form.serialize(),function (res) {
                $("#btnMulti").addClass("disabled");
                if (res.errcode==0){
                    layer.msg("添加成功");
                    setTimeout(function () {
                        location.reload(true);
                    },1500)
                }else{
                    layer.msg(res.message);
                }
                $("#btnCrawl").removeClass("disabled");
            });
        }
    });

    /**
     * 文档目录树
     */
    $("#sidebar").jstree({
        'plugins': ["wholerow", "types", 'dnd', 'contextmenu'],
        "types": {
            "default": {
                "icon": false  // 删除默认图标
            }
        },
        'core': {
            'check_callback': true,
            "multiple": false,
            'animation': 0,
            "data": window.documentCategory
        },
        "contextmenu": {
            show_at_node: false,
            select_node: false,
            "items": {
                "添加文档": {
                    "separator_before": false,
                    "separator_after": true,
                    "_disabled": false,
                    "label": "添加文档",
                    "icon": "fa fa-plus",
                    "action": function (data) {
                        var inst = $.jstree.reference(data.reference),
                            node = inst.get_node(data.reference);

                        openCreateCatalogDialog(node);
                    }
                },
                "删除": {
                    "separator_before": false,
                    "separator_after": true,
                    "_disabled": false,
                    "label": "删除",
                    "icon": "fa fa-trash-o",
                    "action": function (data) {
                        var inst = $.jstree.reference(data.reference);
                        var node = inst.get_node(data.reference);
                        openDeleteDocumentDialog(node);
                    }
                }
            }
        }
    }).on('loaded.jstree', function () {
        window.treeCatalog = $(this).jstree();
    }).on('select_node.jstree', function (node, selected, event) {
        if($("#markdown-save").hasClass('change')) {
            if(confirm("编辑内容未保存，需要保存吗？")){
                saveDocument(false,function () {
                    loadDocument(selected);
                });
                return true;
            }
        }
        loadDocument(selected);

    }).on("move_node.jstree",jstree_save);

    $("#documentTemplateModal").on("click",".section>a[data-type]",function () {
        var $this = $(this).attr("data-type");
        var body = $("#template-" + $this).html();
        if (body) {
            window.isLoad = true;
            window.editor.clear();
            window.editor.insertValue(body);
            window.editor.setCursor({line: 0, ch: 0});
            resetEditorChanged(true);
        }
        $("#documentTemplateModal").modal('hide');
    });

    /*
    *   选中节点，id表示文档的id或者identify
    * */
    function selectedNodeById(id){
        $.each(window.documentCategory,function () {
            if(id == this.id || id==this.identify) {
                window.treeCatalog.deselect_all();
                window.treeCatalog.select_node(this);
                return false;
            }
        });
    }
    function selectedNode(node){
        window.treeCatalog.deselect_all();
        window.treeCatalog.select_node(node);
    }

    window.onpopstate=function(e){
        window.onpop = true;
        selectedNode(e.state)
    }

});

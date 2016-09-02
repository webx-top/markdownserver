function setCookie(c_name,value,expiredays){
    var exdate=new Date()
    exdate.setDate(exdate.getDate()+expiredays)
    document.cookie=c_name+ "=" +escape(value)+((expiredays==null) ? "" : ";expires="+exdate.toGMTString());
}
function getCookie(c_name){
    if(document.cookie.length>0){
        c_start=document.cookie.indexOf(c_name + "=")
        if (c_start!=-1) { 
            c_start=c_start + c_name.length+1 
            c_end=document.cookie.indexOf(";",c_start)
            if (c_end==-1) c_end=document.cookie.length
            return unescape(document.cookie.substring(c_start,c_end))
        } 
    }
    return ""
}
function clickInnerlink(evt){
    evt.preventDefault();
    if($('.m-manual').hasClass('manual-mobile-show-left')){
        $('.m-manual').removeClass('manual-mobile-show-left');
    }
    var that=$(evt.target);
    $('.catalog-list').find('a.current').removeClass('current');
    $.get(that.attr('data-url'),{},function(r){
        setCookie('lastVisited',that.attr('data-url'),365);
        $('.view-body').html(r);
        $('.view-body a[href*="//"]').each(function(){
            $(this).attr('target','_blank');
        });
        $('.view-body code[class^="language-"]').each(function(){
            $(this).attr('class',$(this).attr('class').replace(/^language\-/,'prettyprint lang-'));
        });
        $('.view-body a[href$=".md"]').each(function(){
            $(this).click(clickInbodylink);
        });
        
        var scrollDiv=$('.manual-right');
        scrollDiv.animate({scrollTop: 0}, 0);
        that.addClass('current');
        
        $.getScript(window.JS_PATH+'codeHighlight/loader/run_prettify.js?skin=sunburst');
    },'html');
}

function clickInbodylink(evt){
    evt.preventDefault();
    var url=$(evt.target).attr('href');
    var nav=$('.catalog-list').find('a[data-url="'+url+'"]');
    if(nav.length>0){
        $('.catalog-list').find('a.current').removeClass('current');
        nav.addClass('current');
        var scrollDiv=$('.manual-left .manual-catalog');
        scrollDiv.animate({scrollTop: scrollDiv.scrollTop()+nav.offset().top-100}, 500);
    }
    $.get(url,{},function(r){
        setCookie('lastVisited',url,365);
        $('.view-body').html(r);
        $('.view-body a[href*="//"]').each(function(){
            $(this).attr('target','_blank');
        });
        $('.view-body code[class^="language-"]').each(function(){
            $(this).attr('class',$(this).attr('class').replace(/^language\-/,'prettyprint lang-'));
        });
        $('.view-body a[href$=".md"]').each(function(){
            $(this).click(clickInbodylink);
        });
        
        var scrollDiv=$('.manual-right');
        scrollDiv.animate({scrollTop: 0}, 0);

        $.getScript(window.JS_PATH+'codeHighlight/loader/run_prettify.js?skin=sunburst');
    },'html');
}
$(function(){
window.JS_PATH=$('script[src$="/init.js"]:first').attr('src').split('/init.js')[0]+'/';
var url=$('.catalog-list').data('url');
if(!url)url='SUMMARY.md';
$('.catalog-list').load(url,function(){
    $('.catalog-list a').each(function(){
        var url=$(this).attr('href');
        $(this).attr('href','?'+url);
        $(this).attr('data-url',url);
    });
    $('.catalog-list a').click(clickInnerlink);
    if(window.location.search){
        var url=window.location.search.replace(/'"/g,'')+window.location.hash;
        var a=$('.catalog-list a[href="'+url+'"]');
        if(a.length>0){
            a.trigger('click');
            
            var scrollDiv=$('.manual-left .manual-catalog');
            scrollDiv.animate({scrollTop: scrollDiv.scrollTop()+a.offset().top-100}, 500);
            return;
        }
    }else{
        var lastVisitedURL=getCookie('lastVisited');
        if(lastVisitedURL&&$('.catalog-list a[data-url="'+lastVisitedURL+'"]').length>0){
            $('.catalog-list a[data-url="'+lastVisitedURL+'"]').trigger('click');
            $('body').prepend('<div id="goto-last-visited-tips" style="position:absolute;right:0;border-radius:0 0 0 3px;padding:10px 20px;background:#2d3143;color:white;z-index:999">已经自动切换到您上次浏览的页面</div>');
            window.setTimeout(function(){
                $('#goto-last-visited-tips').fadeOut(500);
            },5000);
            $('#goto-last-visited-tips').click(function(){
                $(this).hide();
            });
            return;
        }
    }
    $('.catalog-list a:first').trigger('click');
});
$('.manual-head .slidebar .icon-menu').click(function(){
    if($('.m-manual').hasClass('manual-mobile-show-left')){
        $('.m-manual').removeClass('manual-mobile-show-left');
    }else{
        $('.m-manual').addClass('manual-mobile-show-left');
    }
});
});
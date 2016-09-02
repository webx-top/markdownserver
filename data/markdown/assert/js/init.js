function clickInnerlink(evt){
    evt.preventDefault();
    if($('.m-manual').hasClass('manual-mobile-show-left')){
        $('.m-manual').removeClass('manual-mobile-show-left');
    }
    var that=$(evt.target);
    $('.catalog-list').find('a.current').removeClass('current');
    $.get(that.attr('data-url'),{},function(r){
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
$('.catalog-list').load('SUMMARY.md',function(){
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
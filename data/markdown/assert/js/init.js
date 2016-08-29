function clickInnerlink(evt){
    evt.preventDefault();
    if($('.m-manual').hasClass('manual-mobile-show-left')){
        $('.m-manual').removeClass('manual-mobile-show-left');
    }
    var that=$(evt.target);
    $('.catalog-list').find('a.current').removeClass('current');
    $.get(that.attr('href'),{},function(r){
        $('.view-body').html(r);
        $('.view-body a[href*="//"]').each(function(){
            $(this).attr('target','_blank');
        });
        $('.view-body a[href$=".md"]').each(function(){
            $(this).click(clickInbodylink);
        });
        that.addClass('current');
    },'html');
}

function clickInbodylink(evt){
    evt.preventDefault();
    var url=$(evt.target).attr('href');
    var nav=$('.catalog-list').find('a[href="'+url+'"]');
    if(nav.length>0){
        $('.catalog-list').find('a.current').removeClass('current');
        nav.addClass('current');
        var scrollDiv=$('.manual-left .manual-catalog');
        scrollDiv.animate({scrollTop: scrollDiv.scrollTop()+nav.offset().top-100}, 1000);
    }
    $.get(url,{},function(r){
        $('.view-body').html(r);
        $('.view-body a[href*="//"]').each(function(){
            $(this).attr('target','_blank');
        });
        $('.view-body a[href$=".md"]').each(function(){
            $(this).click(clickInbodylink);
        });
    },'html');
}
$(function(){
$('.catalog-list').load('SUMMARY.md',function(){
    $('.catalog-list a').click(clickInnerlink);
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
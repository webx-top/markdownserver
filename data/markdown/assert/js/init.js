$(function(){
$('.catalog-list').load('SUMMARY.md',function(){
    $('.catalog-list a').click(function(evt){
        evt.preventDefault();
        if($('.m-manual').hasClass('manual-mobile-show-left')){
            $('.m-manual').removeClass('manual-mobile-show-left');
        }
        var that=$(this);
        $('.catalog-list').find('a.current').removeClass('current');
        $.get($(this).attr('href'),{},function(r){
            $('.view-body').html(r);
            that.addClass('current');
        },'html');
    });
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
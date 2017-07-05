function sync(el, url, type) {
  el.disabled = true;

  fetch(url, {
    method: 'post'
  }).then(function (resp) {
    el.disabled = false;
    alert('sync ' + type + ' success');

  }).catch(function (err) {
    el.disabled = false;
    console.log(err);
    alert('sync ' + type + ' failed');
  });
}

var tickerButton = document.querySelector('.btn-sync-ticker');
if(tickerButton) {
  tickerButton.onclick = function(evt) {
    sync(this, '/sync/ticker', 'ticker');
  }
}

var exchangeButton = document.querySelector('.btn-sync-exchange');
if(exchangeButton) {
  exchangeButton.onclick = function (evt) {
    sync(this, '/sync/exchange', 'exchange');
  }
}

var lendingButton = document.querySelector('.btn-sync-lending');
if(lendingButton) {
  lendingButton.onclick = function (evt) {
    sync(this, '/sync/lending', 'lending');
  }
}

var balanceButton = document.querySelector('.btn-sync-balance');
if(balanceButton) {
  balanceButton.onclick = function (evt) {
    sync(this, '/sync/balance', 'balance');
  }
}

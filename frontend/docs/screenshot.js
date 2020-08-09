document
  .querySelectorAll('.wly-table-body > table > tbody th')
  .forEach((item) => {
    if (isNaN(Number(item.textContent))) {
      let s = ''
      for (let i = 0; i < item.textContent.length; i++) {
        s += 'X'
      }
      item.innerHTML = s
    }
  })

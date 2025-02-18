document.addEventListener('DOMContentLoaded', function () {
  const form = document.getElementById('calculatorForm');
  const resultsDiv = document.getElementById('results');

  form.addEventListener('submit', function (e) {
    e.preventDefault();

    const data = {
      power: parseFloat(document.getElementById('power').value) || 0,
      electricity:
        parseFloat(document.getElementById('electricity').value) || 0,
      deviation1: parseFloat(document.getElementById('deviation1').value) || 0,
      deviation2: parseFloat(document.getElementById('deviation2').value) || 0,
    };

    fetch('/calculator', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(data),
    })
      .then(response => {
        if (!response.ok) {
          throw new Error('Помилка сервера');
        }
        return response.json();
      })
      .then(data => {
        displayResults(data);
      })
      .catch(error => {
        alert('Помилка: ' + error.message);
      });
  });

  function displayResults(results) {
    document.getElementById('profitBefore').textContent =
      results.profitBefore.toFixed(2) + ' тис.грн';
    document.getElementById('profitAfter').textContent =
      results.profitAfter.toFixed(2) + ' тис.грн';

    resultsDiv.style.display = 'block';
  }
});

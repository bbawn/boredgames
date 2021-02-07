let model;
let view;
// let controller;

function SetModel() {
  this.game = null;
  this.NewGame = function() {
    return new Promise((resolve, reject) => {
      const url = new URL('/sets', document.location.origin);
      const d = {
        usernames: ['p1']
      }

      fetch(url, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(d)
      }).then((response) => {
        if (! response.ok) {
          console.log('Network request for ' + url + ' failed with response ' +
                      response.status + ': ' + response.statusText);
          reject();
        }
        console.log('response.body', response.body);
        return response.json().then((json) => {
          this.game = json;
          resolve();
        });
      });
    });
  }
}

function renderBoard(board) {
  console.log('renderBoard', board);
  const gridPanelDiv = document.getElementById('left-panel');
  let grid = document.getElementById('grid');
  if (grid) {
    gridPanelDiv.removeChild(grid);
  }
  grid = document.createElement('table');

  grid.id = 'grid';
  grid.className = 'noselect';
  for (let i = 0; i < 3; i++) {
    const row = document.createElement('tr');

    row.id = 'r' + i;
    row.className = 'grid-row';
    grid.appendChild(row);
    for (let j = 0; j < 4; j++) {
      const cell = document.createElement('td');
      const cellDiv = document.createElement('div');
      const img = document.createElement('img');
      img.name = 'i' + i + '-' + j;
      img.src = '/img/' + board[i * 4 + j] + '.gif'
      console.log('i', i, 'j', j, 'board', board[i*4+j]);

      cellDiv.appendChild(img);
      cell.appendChild(cellDiv);
      row.appendChild(cell);

      cell.className = 'grid-cell';
      cell.id = 'c' + i + '-' + j;
    }
  }
  gridPanelDiv.appendChild(grid);
}

function initialize() {
  const newButton = document.getElementById('new');
  const helpButton = document.getElementById('help');

  model = new SetModel();
  newButton.onclick = function() {
    model.NewGame()
    .then(() => {
      console.log('model.game', model.game);
      // XXX case problem board v Board
      renderBoard(model.game.Board);
    });
  };
  helpButton.onclick = function() {
    /* location.assign('help.html'); */
    window.open('help.html', '_blank');
  };
}

document.addEventListener('DOMContentLoaded', function(event) {
  initialize();
});


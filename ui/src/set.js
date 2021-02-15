let model;

function call(method, path, data) {
  return new Promise((resolve, reject) => {
    const url = new URL(path, document.location.origin);
    fetch(url, {
      method: method,
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(data)
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

// Model for the Set Game
function SetModel() {
  this.game = null;

  // Call the given API method and update game with the result
  this.Call = function(method, path, data) {
    return new Promise((resolve, reject) => {
      const url = new URL(path, document.location.origin);
      fetch(url, {
        method: method,
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(data)
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

  // Create a new game
  this.NewGame = function() {
    const d = { usernames: ['p1'] }
    return this.Call('POST', '/sets', d);
  }

  // Claim a set from current board
  this.ClaimSet = function(set) {
    const d = { username: 'p1', cards: [ set[0], set[1], set[2] ] };
    return this.Call('POST', '/sets/' + this.game.ID + '/claim', d);
  }
}

function renderBoard(board) {
  console.log('renderBoard', board);
  const main = document.getElementById('main');
  let grid = document.getElementById('grid');
  if (grid) {
    main.removeChild(grid);
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
      const a = document.createElement('a');
      const img = document.createElement('img');
      const card = board[i * 4 + j]
      img.alt = card;
      img.src = '/img/' + card + '.gif'
      img.onclick = toggleCellSelected;

      console.log('i', i, 'j', j, 'board', board[i*4+j]);

      a.appendChild(img);
      cell.appendChild(a);
      row.appendChild(cell);

      cell.className = 'grid-cell';
      cell.id = card;
    }
  }
  main.appendChild(grid);
}

function toggleCellSelected(ev) {
  console.log('toggleCellSelected: this ', this);
  console.log('toggleCellSelected: ev', ev);
  const td = this.parentElement.parentElement;
  if (td.classList.contains('grid-selected')) {
    console.log(' toggle remove');
    td.classList.remove('grid-selected');
  } else {
    console.log(' toggle add');
    td.classList.add('grid-selected');
  }
  checkGridForSet()
}

function checkGridForSet() {
  var selectedTds = document.getElementsByClassName('grid-selected');
  if (selectedTds.length > 2) {
    console.log('checkGridForSet selectedTds[0]', selectedTds[0]);
    // selectedTds.foreach(td => set.push(td));
    const set = Array.from(selectedTds).map(td => td.id);
    model.ClaimSet(set);
  }
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
  console.log('DOMContentLoaded', event);
  initialize();
});


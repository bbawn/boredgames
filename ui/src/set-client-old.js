
function assert(condition, message) {
  if (!condition) {
    throw new Error(message || 'Assertion failed');
  }
}

function removeAllChildren(node) {
  while (node.firstChild) {
    node.removeChild(node.firstChild);
  }
}

/* Return random permutation of array using Fisher-Yates */
function shuffleArray(array) {
  for (var i = array.length - 1; i > 0; i--) {
    var j = Math.floor(Math.random() * (i + 1));
    var temp = array[i];
    array[i] = array[j];
    array[j] = temp;
  }
  return array;
}

function makeUserWordLi(word, valid) {
  var li = document.createElement('li');

  li.className = 'word-item list-group-item';
  li.appendChild(makeWordAnchor(word, !valid));

  return li;
}

function makeAnswerWordLi(word, count) {
  var li = document.createElement('li');
  var displayWord = (count > 1 ? word + '(' + count + ')' : word);

  li.className = 'word-item list-group-item';
  li.appendChild(makeWordAnchor(displayWord, false));

  return li;
}

function makeWordEntryLi(word, inputId, buttonText, placeholder, onButtonClick) {
  var li = document.createElement('li');
  var divOuter = document.createElement('div');
  var divAddon = document.createElement('div');
  var input = document.createElement('input');

  input.value = word;
  input.type = 'text';
  input.className = 'form-control word-entry';
  input.id = inputId;
  if (placeholder) {
    input.placeholder = placeholder;
  }
  divAddon.className = "word-entry-addon input-group-addon";
  divAddon.textContent = buttonText;
  divAddon.onclick = onButtonClick;
  divOuter.className = "input-group mb-2 mb-sm-0";
  li.className = "word-entry-item list-group-item";

  divOuter.appendChild(divAddon);
  divOuter.appendChild(input);
  li.appendChild(divOuter);

  return li;
}

function makeGrid(rows, columns, board) {
  var grid = document.createElement('table');

  grid.id = 'grid';
  grid.className = 'noselect';
  for (var i = 0; i < rows; i++) {
    var row = document.createElement('tr');

    row.id = 'r' + i;
    row.className = 'grid-row';
    grid.appendChild(row);
    for (var j = 0; j < columns; j++) {
      var cell = document.createElement('td');
      var cellDiv = document.createElement('div');

      cell.appendChild(cellDiv);
      row.appendChild(cell);

      cell.className = 'grid-cell';
      cell.id = 'c' + i + '-' + j;
      cellDiv.className = 'grid-content';
      cellDiv.textContent = cell.id
    }
  }

  return grid;
}

function touchStart(ev) {
  var elt = document.elementFromPoint(ev.touches[0].clientX, ev.touches[0].clientY);
  pointerStart(ev, elt);
}

function mouseStart(ev) {
  var elt = document.elementFromPoint(ev.clientX, ev.clientY);
  pointerStart(ev, elt);
}

function pointerStart(ev, elt) {
  var grid = document.getElementById('grid');
  if (grid.contains(elt)) {

    /* Disable text selection, touch scrolling. Careful, we need
     * to do this anywhere in grid, not just in grid-content
     * hotspots.
     */
    ev.preventDefault();
  }
  if (! ev.ctrlKey) {
    addUserWord();
  }

  if (elt.classList.contains('grid-content') && 
      (grid.getAttribute('enabled') === '')) {
    addGridCellLetter(elt.parentNode);
  }
}

function touchMove(ev) {
  var elt = document.elementFromPoint(ev.touches[0].clientX, ev.touches[0].clientY);
  var grid = document.getElementById('grid');

  if (grid.contains(elt)) {

    /* Disable text selection, touch scrolling. Careful, we need
     * to do this anywhere in grid, not just in grid-content
     * hotspots.
     */
    ev.preventDefault();
  }

  /* Empirically, elt can be null, not sure how... */
  if (elt && elt.classList.contains('grid-content') && 
      (grid.getAttribute('enabled') === '')) {
    addGridCellLetter(elt.parentNode);
  }
}

function mouseMove(ev) {
  if (ev.buttons == 1) {
    var elt = document.elementFromPoint(ev.clientX, ev.clientY);
    var grid = document.getElementById('grid');

    if (grid.contains(elt)) {

      /* Disable text selection, touch scrolling. Careful, we need
       * to do this anywhere in grid, not just in grid-content
       * hotspots.
       */
      ev.preventDefault();
    }
    if (elt && elt.classList.contains('grid-content') && 
        (grid.getAttribute('enabled') === '')) {

      /* Disable text selection */
      ev.preventDefault();
      addGridCellLetter(elt.parentNode);
    }
  }
}

/* XXX Hackish? consider Node property, e.g disabled */
function setGridState(visible, enabled) {
  var gridContents = document.getElementsByClassName('grid-content');
  var grid = document.getElementById('grid');

  if (enabled) {
    grid.setAttribute('enabled', '');
  } else {
    grid.removeAttribute('enabled');
  }

  for (var i = 0; i < gridContents.length; i++) {
    if (!visible) {
      gridContents[i].classList.add('invisible');
    } else {
      gridContents[i].classList.remove('invisible');
    }
  }
}

function wordEntryKeypress(event) {
  if (event.keyCode == 13) {
    addUserWord();
  }
}

function startKeypress(event) {
  if (event.keyCode == 13) {
    game.changeState(new PlayingState(game));
  }
}

/* Test */
function testAddUserWord(nwords, nchar) {
  var wordEntryTop = document.getElementById('word-entry-top');
  for (var i = 0; i < nwords; i++) {
    wordEntryTop.value = (i).toString()[0].repeat(nchar);
    addUserWord();
  }    
}

function addUserWord() {
  var wordEntryTop = document.getElementById('word-entry-top');
  var ul = document.getElementById('user-word-list');

  if (!wordEntryTop || wordEntryTop.value === '') {
    return;
  }

  var li = makeWordEntryLi(wordEntryTop.value, '', '-', null,
    function(ev) { ul.removeChild(li); });

  // New li goes after top one, simulate insertAfter
  ul.insertBefore(li, wordEntryTop.parentNode.parentNode.nextSibling);
  wordEntryTop.value = '';
  resetWordSelection();
}

function wordEntryButtonClick() {
  addUserWord();
  wordEntryTop.focus();
}

function resetWordSelection() {
  var gridCells = document.getElementsByClassName('grid-cell');

  for (var i = 0; i < gridCells.length; i++) {
    gridCells[i].removeAttribute('selected');
  }
}

function addGridCellLetter(gridCell) {
  var wordEntryTop = document.getElementById('word-entry-top');
  if (wordEntryTop && (gridCell.getAttribute('selected') !== '')) {
    wordEntryTop.value += gridCell.textContent.toLowerCase();
    gridCell.setAttribute('selected', '');
  }
}

/* Use state design pattern to manage state machine for the game.
 * See http://www.dofactory.com/javascript/state-design-pattern
 * The state machine is:
 *
 * State        Transitions to
 * -----        --------------
 * Start        Playing
 * Playing      Paused, SetClaimed, Done
 * Paused       Playing
 * SetClaimed   Playing
 * Done         Start
 */

var GameContext = function() {
  this.rows = 3;             // XXX move to config object?
  this.columns = 4;          // XXX move to config object?
    this.game = null;
  this.currentState = new StartState(this);

  this.changeState = function(state) {
    this.currentState = state;
    this.currentState.go();
  };

  this.start = function() {
    this.currentState.go();
  };
};

var StartState = function(game) {
  this.game = game;

  this.go = function() {
    var gridPanelDiv = document.getElementById('left-panel');
    var playingButton = document.getElementById('playing');
    var startButton = document.getElementById('start-button');
    var yourWordsLabel = document.getElementById('your-words-label');
    var userWordList = document.getElementById('user-word-list');
    var grid = document.getElementById('grid');
    var solutionLabel = document.getElementById('solution-label');
    var solutionWordList = document.getElementById('solution-word-list');
    var rightPanel = document.getElementById('right-panel');
    var startModal = document.getElementById('start-modal');

    game.userWords = [];
    game.userWords = {};
    game.userScore = undefined;
    game.solutionScore = undefined;
    startButton.onclick = function() {
      game.changeState(new PlayingState(game));
    };
    startButton.onkeypress = startKeypress;

    /* Helpful to start with focus in word entry for keyboard entry
     * but it's super annoying on phones - pops up the keyboard when
     * user probably should use touch entry
     */
    if (isTouchDevice()) {
      startButton.focus();
    }
    playingButton.innerHTML = 'Play';
    if (grid) {
      gridPanelDiv.removeChild(grid);
    }
    gridPanelDiv.appendChild(makeGrid(game.rank, game.board));
    setGridState(false, false);
    playingButton.disabled = false;
    yourWordsLabel.textContent = 'Your words:';
    removeAllChildren(userWordList);

    var li = makeWordEntryLi('', 'word-entry-top', '+', 'Enter word', 
        wordEntryButtonClick);
    userWordList.appendChild(li);
    var wordEntryInput = document.getElementById('word-entry-top');
    wordEntryInput.disabled = true;
    wordEntryInput.onkeypress = wordEntryKeypress;

    solutionLabel.innerHTML = '';
    removeAllChildren(solutionWordList);
    rightPanel.style.display = 'none';
    startModal.style.display = 'block';
  };
};

var PlayingState = function(game) {
  this.game = game;

  this.go = function() {
    var playingButton = document.getElementById('playing');
    var wordEntry = document.getElementById('word-entry-top');
    var startModal = document.getElementById('start-modal');

    playingButton.onclick = function() {
      game.changeState(new PausedState(game));
    };
    playingButton.innerHTML = 'Pause';
    wordEntry.disabled = false;
    wordEntry.focus();

    setGridState(true, true);
    startModal.style.display = 'none';
  };
};

var PausedState = function(game) {
  this.game = game;

  this.go = function() {
    var playingButton = document.getElementById('playing');
    playingButton.onclick = function() {
      game.changeState(new PlayingState(game));
    };
    playingButton.innerHTML = 'Play';
    document.getElementById('word-entry-top').disabled = true;
    document.getElementById('grid').disabled = true;
    setGridState(false, false);
  };
};

var SolvingState = function(game) {
  this.game = game;

  this.go = function() {
    var wordEntries = document.getElementsByClassName('word-entry');
    var statusHeader = document.getElementById('status');

    statusHeader.textContent = 'Solving...';
    resetWordSelection();

    game.userWords = [];
    for (var i = 0; i < wordEntries.length; i++) {
      if (wordEntries[i].value !== '') {
        game.userWords.push(wordEntries[i].value.toLowerCase());
      }
    }

    // XXX is this the best way to construct the URL?
    var url = new URL('/api', document.location.origin);

    url.searchParams.append('words', game.userWords.join(' '));
    url.searchParams.append('board', game.board.join(' '));

    fetch(url).then(function(response) {

      // State changed during fetch, drop response
      if (game.currentState.constuctor === SolvingState) {
        return;
      }

      // TODO - should alert() or something user-visible...
      if (! response.ok) {
        console.log('Network request for ' + url + 'failed with response ' +
                    response.status + ': ' + response.statusText);
      }
      response.json().then(function(json) {
        game.solutionWords = json.answers; // XXX change server
        game.userScore = json.word_score;
        game.solutionScore = json.answer_score;
        game.changeState(new SolvedState(game));
      });
    });
  };
};

var SolvedState = function(game) {
  this.game = game;

  this.go = function() {
    var playingButton = document.getElementById('playing');
    var yourWordsLabel = document.getElementById('your-words-label');
    var userWordList = document.getElementById('user-word-list');
    var statusHeader = document.getElementById('status');
    var rightPanel = document.getElementById('right-panel');

    // Header
    playingButton.innerHTML = 'Play';
    setGridState(true, false);
    playingButton.disabled = true;

    // User word panel TODO: sort on server??
    yourWordsLabel.textContent = 'Your words (' +
      game.userWords.length + '):';

    validatedWords = validateWords(game.solutionWords, game.userWords);
    removeAllChildren(userWordList);
    var sortedWords = game.userWords.sort();
    for (var i = 0; i < sortedWords.length; i++) {
      var li = makeUserWordLi(sortedWords[i],
                              validatedWords[sortedWords[i]]);
      userWordList.appendChild(li);
    };
    statusHeader.textContent = 'Your score: ' + game.userScore +
                               ' out of ' + game.solutionScore;

    var answerLabel = document.getElementById('solution-label');
    var answerWordList = document.getElementById('solution-word-list');
    answerLabel.textContent = 'Solution (' +
      Object.getOwnPropertyNames(game.solutionWords).length +
      '):';

    // Answer word panel TODO: sort on server??
    var sortedWords = Object.keys(game.solutionWords).sort();
    for (var i = 0; i < sortedWords.length; i++) {
      var li = makeAnswerWordLi(sortedWords[i],
                                game.solutionWords[sortedWords[i]]);
      answerWordList.appendChild(li);
    }

    rightPanel.style.display = '';
  };
};

function toggleFullScreen() {
  console.log('toggleFullScreen');
  if ((document.fullScreenElement && document.fullScreenElement !== null) ||    
   (!document.mozFullScreen && !document.webkitIsFullScreen)) {
    if (document.documentElement.requestFullScreen) {  
      document.documentElement.requestFullScreen();  
    } else if (document.documentElement.mozRequestFullScreen) {  
      document.documentElement.mozRequestFullScreen();  
    } else if (document.documentElement.webkitRequestFullScreen) {  
      document.documentElement.webkitRequestFullScreen(Element.ALLOW_KEYBOARD_INPUT);  
    }  
  } else {  
    if (document.cancelFullScreen) {  
      document.cancelFullScreen();  
    } else if (document.mozCancelFullScreen) {  
      document.mozCancelFullScreen();  
    } else if (document.webkitCancelFullScreen) {  
      document.webkitCancelFullScreen();  
    }  
  }  
}

function isTouchDevice() {
  return 'ontouchstart' in document.documentElement;
}

/* Test */
function testDisplayUserWord(nwords, nchar) {
  var userWordList = document.getElementById('user-word-list');
  for (var i = 0; i < nwords; i++) {
    var li = makeUserWordLi((i).toString()[0].repeat(nchar), nchar % 2);
    userWordList.appendChild(li);
  };
}
function testDisplayAnswerWord(nwords, nchar) {
  var answerWordList = document.getElementById('solution-word-list');
  for (var i = 0; i < nwords; i++) {
    var li = makeAnswerWordLi((i).toString()[0].repeat(nchar), nchar % 2);
    answerWordList.appendChild(li);
  }
}


var game;
function initialize() {
  var newButton = document.getElementById('new');
  var helpButton = document.getElementById('help');
  game = new GameContext();

  newButton.onclick = function() {
    game.changeState(new StartState(game));
  };
  helpButton.onclick = function() {
    /* location.assign('help.html'); */
    window.open('help.html', '_blank');
  };

  document.addEventListener('touchstart', touchStart, {passive: false});
  document.addEventListener('mousedown', mouseStart, {passive: false});

  /* touchenter is, apparently, not (yet) implemented */
  document.addEventListener('touchmove', touchMove, {passive: false});
  document.addEventListener('mousemove', mouseMove, {passive: false});

  game.start();
}

document.addEventListener('DOMContentLoaded', function(event) {
  initialize();
});


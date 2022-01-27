let model;

// Model for the Set Game
function SetModel() {
  this.config = { machineHandicap: 60 };
  this.localUsername = "p1";
  this.game = null;

  // Call the given API method and update game with the result
  this.Call = function (method, path, data) {
    return new Promise((resolve, reject) => {
      const url = new URL(path, document.location.origin);
      fetch(url, {
        method: method,
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(data),
      }).then((response) => {
        if (!response.ok) {
          console.log(
            "Network request for " +
              url +
              " failed with response " +
              response.status +
              ": " +
              response.statusText
          );
          reject();
        }
        console.log("response.body", response.body);
        return response.json().then((json) => {
          this.game = json;
          resolve();
        });
      });
    });
  };

  // Create a new game
  this.NewGame = function () {
    const d = { usernames: ["p1", "machine"] };
    return this.Call("POST", "/sets", d);
  };

  // Claim a set from current board
  this.ClaimSet = function (username, set) {
    const d = { username: username, cards: [set[0], set[1], set[2]] };
    return this.Call("POST", "/sets/" + this.game.ID + "/claim", d);
  };

  // Expand the board
  this.Expand = function () {
    return this.Call("POST", "/sets/" + this.game.ID + "/expand");
  };

  // Start the next round
  this.Next = function () {
    return this.Call("POST", "/sets/" + this.game.ID + "/next");
  };
}

// getState returns the state of the game
function getState(game) {
  if (game.ClaimedUsername === "") {
    return "Playing";
  } else {
    return "SetClaimed";
  }
}

function render(game) {
  console.log("render", game);

  const main = document.getElementById("main");
  let board = document.getElementById("board");
  if (board) {
    main.removeChild(board);
  }
  board = createBoard(game);
  main.appendChild(board);

  let scorecard = document.getElementById("scorecard");
  if (scorecard) {
    main.removeChild(scorecard);
  }
  scorecard = createScorecard(game);
  main.appendChild(scorecard);

  setMessage(game);
}

function setMessage(game) {
  const message = document.getElementById("message");
  removeAllChildren(message);
  const state = getState(game);
  if (state == "Playing") {
    // TODO: compare old game and see if any player's score decreased and add
    // to message
    let text = document.createTextNode("Select a Set");
    message.appendChild(text);
  } else if (state == "SetClaimed") {
    const span = document.createElement("span");
    let text = document.createTextNode(game.ClaimedUsername + " claimed ");
    span.appendChild(text);
    message.appendChild(span);
    for (const card of game.ClaimedSet) {
      const img = document.createElement("img");
      img.className = "message-img";
      img.alt = card;
      img.src = "/img/" + card + ".gif";
      message.appendChild(img);
    }
  }
}

function createBoard(game) {
  const board = document.createElement("table");

  board.id = "board";
  for (let i = 0; i < 3; i++) {
    const row = document.createElement("tr");

    row.id = "r" + i;
    row.className = "board-row";
    board.appendChild(row);
    const nCols = game.Board.length / 3;
    for (let j = 0; j < nCols; j++) {
      const cell = document.createElement("td");
      const a = document.createElement("a");
      const img = document.createElement("img");
      const card = game.Board[i * nCols + j] || `empty_card`;
      img.alt = card;
      img.src = "/img/" + card + ".gif";
      cell.className = "board-img";
      if (getState(game) === "Playing") {
        img.onclick = toggleCellSelected;
      }

      a.appendChild(img);
      cell.appendChild(a);
      row.appendChild(cell);

      cell.className = "board-cell";
      cell.id = card;
    }
  }
  return board;
}

function createScorecard(game) {
  const scorecard = document.createElement("table");
  scorecard.id = "scorecard";
  const head = document.createElement("tr");
  let cell = document.createElement("td");
  cell.className = "score-header";
  let cellText = document.createTextNode("Player");
  cell.appendChild(cellText);
  head.appendChild(cell);
  cell = document.createElement("td");
  cell.className = "score-header";
  cellText = document.createTextNode("Handicap");
  cell.appendChild(cellText);
  head.appendChild(cell);
  cell = document.createElement("td");
  cell.className = "score-header";
  cellText = document.createTextNode("Score");
  cell.appendChild(cellText);
  head.appendChild(cell);
  scorecard.appendChild(head);

  console.log("createScorecard: game", game);
  for (const [username, player] of Object.entries(game.Players)) {
    const row = document.createElement("tr");
    let cell = document.createElement("td");
    cell.className = "score-cell";
    let cellText = document.createTextNode(username);
    cell.appendChild(cellText);
    row.appendChild(cell);
    cell = document.createElement("td");
    cell.className = "score-cell";
    if (username.toLowerCase() === "machine") {
      cellInput = document.createElement("input");
      cellInput.type = "number";
      cellInput.value = model.config.machineHandicap;
      cellInput.addEventListener("input", onHandicapInput);
      cellInput.className = "handicap-cell";
      cell.appendChild(cellInput);
    }
    row.appendChild(cell);
    cell = document.createElement("td");
    cell.className = "score-cell";
    cellText = document.createTextNode(player.Sets.length);
    cell.appendChild(cellText);
    row.appendChild(cell);
    scorecard.appendChild(row);
  }
  return scorecard;
}

function onHandicapInput(e) {
  model.config.machineHandicap = parseInt(e.target.value);
  console.log(
    "onHandicapInpute machineHandicap",
    model.config.machineHandicap
  );
}

function toggleCellSelected(ev) {
  console.log("toggleCellSelected: this ", this);
  console.log("toggleCellSelected: ev", ev);
  const td = this.parentElement.parentElement;
  if (td.classList.contains("board-selected")) {
    console.log(" toggle remove");
    td.classList.remove("board-selected");
  } else {
    console.log(" toggle add");
    td.classList.add("board-selected");
  }
  checkBoardForSet();
}

function checkBoardForSet() {
  var selectedTds = document.getElementsByClassName("board-selected");
  if (selectedTds.length > 2) {
    console.log("checkBoardForSet selectedTds[0]", selectedTds[0]);
    const set = Array.from(selectedTds).map((td) => td.id);
    model.ClaimSet(model.localUsername, set).then(() => {
      console.log("claim() response: model.game", model.game);
      render(model.game);
    });
  }
  cancelMachineTimer();
}

function removeAllChildren(parent) {
  while (parent.firstChild) {
    parent.removeChild(parent.firstChild);
  }
}

let machineTimer;
function scheduleMachineTimer() {
  console.log(
    Date.now(),
    "machineClaim: machineHandicap",
    model.config.machineHandicap
  );
  clearTimeout(machineTimer);
  machineTimer = setTimeout(
    machineClaim,
    model.config.machineHandicap * 1000
  );
}

function cancelMachineTimer() {
  clearTimeout(machineTimer);
}

function machineClaim() {
  console.log(Date.now(), "machineClaim");
  const set = findSet(model.game.Board);
  if (!set) {
    console.log("machine failed to find a set");
    return;
  }
  model.ClaimSet("machine", set).then(() => {
    console.log("claim() response: model.game", model.game);
    render(model.game);
  });
}

function findSet(b) {
  // Enumerate each set on board
  for (let i = 0; i < b.length - 2; i++) {
    const c1 = b[i];
    if (c1) {
      for (let j = i + 1; j < b.length - 1; j++) {
        const c2 = b[j];
        if (c2) {
          for (let k = j + 1; k < b.length; k++) {
            const c3 = b[k];
            if (c3) {
              const cs = [c1, c2, c3];
              if (isSet(cs)) {
                return cs;
              }
            }
          }
        }
      }
    }
  }
}

function isSet(cs) {
  if (cs[0] == cs[1] || cs[1] == cs[2] || cs[0] == cs[2]) {
    return false;
  }
  let ret = true;
  for (let i = 0; i < 4; i++) {
    ret &= setMatch(cs[0][i], cs[1][i], cs[2][i]);
  }
  return ret;
}

// setMatch return true if the given values are a "match" by Set rules:
// meaning they are either all the same or all different
function setMatch(b1, b2, b3) {
  return (b1 == b2 && b2 == b3) || (b1 != b2 && b2 != b3 && b1 != b3);
}

function initialize() {
  const newButton = document.getElementById("new");
  const dealButton = document.getElementById("deal");
  const helpButton = document.getElementById("help");

  model = new SetModel();
  newButton.onclick = function () {
    model.NewGame().then(() => {
      console.log("new model.game", model.game);
      scheduleMachineTimer();
      render(model.game);
    });
  };
  dealButton.onclick = function () {
    if (getState(model.game) === "Playing") {
      model.Expand().then(() => {
        console.log("expand() response: model.game", model.game);
        scheduleMachineTimer();
        render(model.game);
      });
    } else {
      model.Next().then(() => {
        console.log("next() response: model.game", model.game);
        scheduleMachineTimer();
        render(model.game);
      });
    }
  };
  helpButton.onclick = function () {
    window.open("help.html", "_blank");
  };
}

document.addEventListener("DOMContentLoaded", function (event) {
  console.log("DOMContentLoaded", event);
  initialize();
});

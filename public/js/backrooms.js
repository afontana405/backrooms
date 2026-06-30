/* Backrooms — Level 0. Procedurally generated liminal space.
 * Mono-yellow chevron wallpaper, worn carpet, drop-ceiling fluorescents + hum.
 * Click to look, WASD to move, Esc to release, M to mute the buzz.
 * Walls only (no enemies/puzzles yet).
 */
(function () {
  if (!window.THREE) return;
  var THREE = window.THREE;
  var CELL = 3.2, WALL_H = 3.0, EYE = 1.55;

  function ri(a, b) { return a + Math.floor(Math.random() * (b - a + 1)); }

  // ── Procedural Map Generation ──
  var MIN_ROOMS = 10, MAX_ROOMS = 12;
  var RM = 9, HALF_RM = Math.floor(RM / 2);
  var ROOM_COUNT = ri(MIN_ROOMS, MAX_ROOMS);

  var GW = Math.ceil(Math.sqrt(ROOM_COUNT) * 14) + 20;
  var GH = Math.ceil(Math.sqrt(ROOM_COUNT) * 14) + 20;

  var g = [];
  for (var r = 0; r < GH; r++) { g[r] = []; for (var c = 0; c < GW; c++) g[r][c] = 1; }

  function carve(x0, y0, w, h) {
    for (var y = y0; y < y0 + h; y++) for (var x = x0; x < x0 + w; x++)
      if (x >= 0 && x < GW && y >= 0 && y < GH) g[y][x] = 0;
  }

  var rooms = [];
  var doors = [];

  function addRoom(cx, cy) {
    carve(cx - HALF_RM, cy - HALF_RM, RM, RM);
    rooms.push({ cx: cx, cy: cy });
  }

  // Place first room near center
  addRoom(Math.floor(GW / 2) + ri(-10, 10), Math.floor(GH / 2) + ri(-10, 10));

  // Generate remaining rooms with hallways
  for (var i = 1; i < ROOM_COUNT; i++) {
    var p = rooms[i-1];
    var ok = false;
    for (var att = 0; att < 300 && !ok; att++) {
      var angle = Math.random() * Math.PI * 2;
      var dist = ri(6, 14) + RM;
      var nx = Math.round(p.cx + Math.cos(angle) * dist + ri(-3, 3));
      var ny = Math.round(p.cy + Math.sin(angle) * dist + ri(-3, 3));

      if (nx - HALF_RM - 1 < 1 || nx + HALF_RM + 1 >= GW - 1) continue;
      if (ny - HALF_RM - 1 < 1 || ny + HALF_RM + 1 >= GH - 1) continue;

      var overlap = false;
      for (var r = 0; r < rooms.length; r++) {
        if (Math.abs(nx - rooms[r].cx) < RM + 3 && Math.abs(ny - rooms[r].cy) < RM + 3) {
          overlap = true; break;
        }
      }
      if (overlap) continue;

      // Pre-compute hallway cells for validation
      var hc = [];
      var bendX = Math.random() < 0.5 ? p.cx : nx;
      var bendY = Math.random() < 0.5 ? ny : p.cy;

      var stepX = p.cx < bendX ? 1 : -1;
      for (var hx = p.cx; hx !== bendX; hx += stepX) hc.push([hx, p.cy]);
      var stepY = p.cy < bendY ? 1 : -1;
      for (var hy = p.cy; hy !== bendY; hy += stepY) hc.push([bendX, hy]);

      var s2x = bendX < nx ? 1 : -1;
      for (var hx = bendX; hx !== nx; hx += s2x) hc.push([hx, bendY]);
      var s2y = bendY < ny ? 1 : -1;
      for (var hy = bendY; hy !== ny; hy += s2y) hc.push([nx, hy]);

      // Validate hallway doesn't intersect any existing walkable area outside source room
      var srcMinX = p.cx - HALF_RM, srcMaxX = p.cx + HALF_RM;
      var srcMinY = p.cy - HALF_RM, srcMaxY = p.cy + HALF_RM;
      var valid = true;
      for (var h = 0; h < hc.length && valid; h++) {
        var hx = hc[h][0], hy = hc[h][1];
        if (hx >= srcMinX && hx <= srcMaxX && hy >= srcMinY && hy <= srcMaxY) continue;
        if (g[hy][hx] === 0) valid = false;
      }
      // Validate candidate room cells are also clear
      for (var ry = ny - HALF_RM; ry <= ny + HALF_RM && valid; ry++) {
        for (var rx = nx - HALF_RM; rx <= nx + HALF_RM && valid; rx++) {
          if (g[ry][rx] === 0) valid = false;
        }
      }
      if (!valid) continue;

      // All clear — carve room
      addRoom(nx, ny);

      // Carve hallway
      for (var h = 0; h < hc.length; h++) g[hc[h][1]][hc[h][0]] = 0;

      // Door: 2 cells at the hallway exit from the source room
      var exitIdx = -1;
      for (var h = 0; h < hc.length; h++) {
        var hx = hc[h][0], hy = hc[h][1];
        if (hx < p.cx - HALF_RM || hx > p.cx + HALF_RM || hy < p.cy - HALF_RM || hy > p.cy + HALF_RM) {
          exitIdx = h; break;
        }
      }
      var dc = [];
      if (exitIdx >= 0 && exitIdx + 1 < hc.length) dc = [hc[exitIdx], hc[exitIdx + 1]];
      else dc = [[Math.max(1,p.cx+HALF_RM+1), p.cy], [Math.max(1,p.cx+HALF_RM+2), p.cy]];

      for (var d = 0; d < dc.length; d++) g[dc[d][1]][dc[d][0]] = 1;
      doors.push({ cells: dc, room: i - 1 });
      ok = true;
    }
  }

  function isWall(c, r) { if (r < 0 || r >= GH || c < 0 || c >= GW) return true; return g[r][c] === 1; }

  var spawn = rooms[0], exitR = rooms[rooms.length - 1];
  var sx = spawn.cx * CELL + CELL / 2, sz = spawn.cy * CELL + CELL / 2;
  var ex = exitR.cx * CELL + CELL / 2, ez = exitR.cy * CELL + CELL / 2;

  // ── Procedural textures ──
  function cv(sz) { var c = document.createElement('canvas'); c.width = c.height = sz; return c; }
  function wallpaperTex() {
    var c = cv(128), x = c.getContext('2d');
    x.fillStyle = '#cabd6e'; x.fillRect(0, 0, 128, 128);
    // vertical chevron stripes
    for (var i = 0; i < 128; i += 16) {
      x.strokeStyle = i % 32 === 0 ? '#b3a352' : '#d6ca7e';
      x.lineWidth = 6; x.beginPath();
      for (var yy = -16; yy < 144; yy += 16) { x.moveTo(i, yy); x.lineTo(i + 8, yy + 8); x.lineTo(i, yy + 16); }
      x.stroke();
    }
    var t = new THREE.CanvasTexture(c); t.wrapS = t.wrapT = THREE.RepeatWrapping; return t;
  }
  function carpetTex() {
    var c = cv(128), x = c.getContext('2d');
    x.fillStyle = '#7d7235'; x.fillRect(0, 0, 128, 128);
    for (var i = 0; i < 2600; i++) {
      var v = ri(-18, 18);
      x.fillStyle = 'rgb(' + (125 + v) + ',' + (114 + v) + ',' + (53 + v) + ')';
      x.fillRect(ri(0, 127), ri(0, 127), 2, 2);
    }
    var t = new THREE.CanvasTexture(c); t.wrapS = t.wrapT = THREE.RepeatWrapping; return t;
  }
  function ceilTex() {
    var c = cv(128), x = c.getContext('2d');
    x.fillStyle = '#d9d0a2'; x.fillRect(0, 0, 128, 128);
    x.strokeStyle = '#9c946c'; x.lineWidth = 4;
    x.strokeRect(0, 0, 64, 64); x.strokeRect(64, 0, 64, 64); x.strokeRect(0, 64, 64, 64); x.strokeRect(64, 64, 64, 64);
    var t = new THREE.CanvasTexture(c); t.wrapS = t.wrapT = THREE.RepeatWrapping; return t;
  }

  // ── Scene ──
  var root = document.getElementById('br-root') || document.body;
  var scene = new THREE.Scene();
  scene.background = new THREE.Color(0x161200);
  scene.fog = new THREE.FogExp2(0x171300, 0.03);

  var camera = new THREE.PerspectiveCamera(72, window.innerWidth / window.innerHeight, 0.1, 800);
  camera.rotation.order = 'YXZ';
  var renderer = new THREE.WebGLRenderer({ antialias: true });
  renderer.setPixelRatio(window.devicePixelRatio || 1);
  renderer.setSize(window.innerWidth, window.innerHeight);
  root.appendChild(renderer.domElement);

  var ambient = new THREE.AmbientLight(0xfff0b4, 0.7); scene.add(ambient);
  scene.add(new THREE.HemisphereLight(0xfff0b4, 0x2e2800, 0.4));

  // Player-tool state + flashlight
  var flashOn = false, battery = 100, stamina = 100, hidden = false, activeHide = null;
  var nearDist = 1e9, bobPhase = 0, stepT = 0, flick = 0;
  var spot = new THREE.SpotLight(0xfff2d0, 0, 26, Math.PI / 7, 0.4, 1.2);
  var spotTarget = new THREE.Object3D(); scene.add(spotTarget); spot.target = spotTarget; scene.add(spot);

  var floorW = GW * CELL, floorD = GH * CELL;
  var ctex = carpetTex(); ctex.repeat.set(GW, GH);
  var floor = new THREE.Mesh(new THREE.PlaneGeometry(floorW, floorD).rotateX(-Math.PI / 2),
    new THREE.MeshStandardMaterial({ map: ctex, roughness: 1 }));
  floor.position.set(floorW / 2, 0, floorD / 2); scene.add(floor);

  var cetex = ceilTex(); cetex.repeat.set(GW, GH);
  var ceil = new THREE.Mesh(new THREE.PlaneGeometry(floorW, floorD).rotateX(Math.PI / 2),
    new THREE.MeshStandardMaterial({ map: cetex, roughness: 1 }));
  ceil.position.set(floorW / 2, WALL_H, floorD / 2); scene.add(ceil);

  // Walls (instanced, wallpaper) — skip doorway cells so opened doors are clear
  var doorSet = {}; doors.forEach(function (d) { d.cells.forEach(function (cl) { doorSet[cl[0] + ',' + cl[1]] = 1; }); });
  var wallCount = 0; for (var wr = 0; wr < GH; wr++) for (var wc = 0; wc < GW; wc++) if (g[wr][wc] === 1 && !doorSet[wc + ',' + wr]) wallCount++;
  var walls = new THREE.InstancedMesh(
    new THREE.BoxGeometry(CELL, WALL_H, CELL),
    new THREE.MeshStandardMaterial({ map: wallpaperTex(), roughness: 0.95 }),
    wallCount);
  var m4 = new THREE.Matrix4(), wi = 0;
  for (var ar = 0; ar < GH; ar++) for (var ac = 0; ac < GW; ac++) {
    if (g[ar][ac] !== 1 || doorSet[ac + ',' + ar]) continue;
    m4.makeTranslation(ac * CELL + CELL / 2, WALL_H / 2, ar * CELL + CELL / 2);
    walls.setMatrixAt(wi++, m4);
  }
  walls.instanceMatrix.needsUpdate = true; scene.add(walls);

  // Fluorescent ceiling panels + a few real lights
  var panelMat = new THREE.MeshBasicMaterial({ color: 0xfff6cc });
  var litThisMany = 0;
  for (var lr = 3; lr < GH; lr += 6) for (var lc = 3; lc < GW; lc += 6) {
    if (g[lr][lc] === 1) continue;
    var panel = new THREE.Mesh(new THREE.PlaneGeometry(CELL * 0.7, CELL * 0.22).rotateX(Math.PI / 2), panelMat);
    panel.position.set(lc * CELL + CELL / 2, WALL_H - 0.03, lr * CELL + CELL / 2);
    scene.add(panel);
    if (litThisMany < 14) {
      var pl = new THREE.PointLight(0xfff0c0, 0.7, CELL * 9, 2);
      pl.position.set(lc * CELL + CELL / 2, WALL_H - 0.4, lr * CELL + CELL / 2);
      scene.add(pl); litThisMany++;
    }
  }

  // ── Sparse odd props ──
  function prop(kind, wx, wz) {
    var grp = new THREE.Group(); grp.position.set(wx, 0, wz);
    var box = function (w, h, d, col, y) { var mm = new THREE.Mesh(new THREE.BoxGeometry(w, h, d), new THREE.MeshStandardMaterial({ color: col, roughness: 0.9 })); mm.position.y = y; return mm; };
    var cyl = function (r, h, col, y) { var mm = new THREE.Mesh(new THREE.CylinderGeometry(r, r, h, 12), new THREE.MeshStandardMaterial({ color: col, roughness: 0.9 })); mm.position.y = y; return mm; };
    if (kind === 0) { grp.add(box(0.5, 0.05, 0.5, 0x3a4a66, 0.5)); grp.add(box(0.05, 0.5, 0.05, 0x2b3850, 0.25)); grp.add(box(0.05, 0.5, 0.05, 0x2b3850, 0.25).translateX(0.22)); grp.add(box(0.5, 0.5, 0.05, 0x2b3850, 0.75).translateZ(-0.22)); } // chair
    else if (kind === 1) { grp.add(box(0.7, 0.7, 0.7, 0xb08a4a, 0.35)); } // cardboard box
    else if (kind === 2) { var pp = cyl(0.12, 2.4, 0x7d7d82, 1.2); pp.rotation.z = Math.PI / 2; grp.add(pp); } // pipe
    else if (kind === 3) { var cn = cyl(0.02, 0.6, 0xff7a18, 0.3); cn.scale.set(1, 1, 1); var base = box(0.5, 0.05, 0.5, 0xff7a18, 0.02); var cone = new THREE.Mesh(new THREE.ConeGeometry(0.22, 0.55, 14), new THREE.MeshStandardMaterial({ color: 0xff7a18 })); cone.position.y = 0.3; grp.add(base); grp.add(cone); } // wet-floor cone
    else { var sign = box(0.9, 0.4, 0.08, 0x0a0a0a, WALL_H - 0.6); var glow = new THREE.Mesh(new THREE.PlaneGeometry(0.8, 0.3), new THREE.MeshBasicMaterial({ color: 0x33ff55 })); glow.position.set(0, WALL_H - 0.6, 0.05); grp.add(sign); grp.add(glow); } // EXIT sign
    scene.add(grp);
  }
  var placed = 0;
  for (var pp2 = 0; pp2 < 400 && placed < 16; pp2++) {
    var cc = ri(2, GW - 3), rr2 = ri(2, GH - 3);
    if (g[rr2][cc] !== 0) continue;
    if (Math.abs(cc - spawn.cx) + Math.abs(rr2 - spawn.cy) < 4) continue;
    prop(ri(0, 3), cc * CELL + CELL / 2, rr2 * CELL + CELL / 2); placed++;
  }
  prop(4, ex, ez); // EXIT sign at the exit room

  // ── Player ──
  var px = sx, pz = sz, yaw = 0, pitch = 0, SPEED = 5.0, PR = 0.42;
  function blocked(x, z) {
    for (var ox = -PR; ox <= PR + 0.001; ox += PR)
      for (var oz = -PR; oz <= PR + 0.001; oz += PR)
        if (isWall(Math.floor((x + ox) / CELL), Math.floor((z + oz) / CELL))) return true;
    return false;
  }

  // ── Entities (pass 1: chaser, still-life, smiler) ──
  function roomCenter(rm) { return { x: rm.cx * CELL + CELL / 2, z: rm.cy * CELL + CELL / 2 }; }
  var checkpoint = { x: sx, z: sz }, lastRoomIdx = -1, inv = 0;

  function humanoid(col) {
    var grp = new THREE.Group();
    var mat = new THREE.MeshStandardMaterial({ color: col, roughness: 0.9 });
    function b(w, h, d, x, y, z) { var mm = new THREE.Mesh(new THREE.BoxGeometry(w, h, d), mat); mm.position.set(x, y, z); grp.add(mm); }
    b(0.5, 0.85, 0.28, 0, 1.15, 0);   // torso
    b(0.3, 0.3, 0.3, 0, 1.75, 0);     // head
    b(0.15, 0.85, 0.15, -0.16, 0.42, 0); b(0.15, 0.85, 0.15, 0.16, 0.42, 0); // legs
    b(0.12, 0.7, 0.12, -0.34, 1.2, 0); b(0.12, 0.7, 0.12, 0.34, 1.2, 0);     // arms
    return grp;
  }
  function smilerModel() {
    var grp = new THREE.Group();
    var body = new THREE.Mesh(new THREE.BoxGeometry(0.7, 2.2, 0.4), new THREE.MeshStandardMaterial({ color: 0x050505, roughness: 1 }));
    body.position.y = 1.1; grp.add(body);
    var eyeMat = new THREE.MeshBasicMaterial({ color: 0xfff36b });
    var e1 = new THREE.Mesh(new THREE.SphereGeometry(0.1, 10, 10), eyeMat); e1.position.set(-0.14, 1.85, 0.21); grp.add(e1);
    var e2 = e1.clone(); e2.position.x = 0.14; grp.add(e2);
    var teethMat = new THREE.MeshBasicMaterial({ color: 0xffffff });
    for (var i = -3; i <= 3; i++) { var t = new THREE.Mesh(new THREE.BoxGeometry(0.05, 0.13, 0.02), teethMat); t.position.set(i * 0.065, 1.62 + Math.abs(i) * 0.02, 0.21); grp.add(t); }
    return grp;
  }
  function hazardModel(col) {
    var grp = new THREE.Group();
    var disc = new THREE.Mesh(new THREE.CircleGeometry(CELL * 0.9, 20).rotateX(-Math.PI / 2), new THREE.MeshBasicMaterial({ color: col, transparent: true, opacity: 0.4 }));
    disc.position.y = 0.04; grp.add(disc);
    for (var i = 0; i < 6; i++) { var lump = new THREE.Mesh(new THREE.SphereGeometry(0.25 + Math.random() * 0.3, 8, 8), new THREE.MeshStandardMaterial({ color: col, roughness: 0.85 })); lump.position.set((Math.random() - 0.5) * CELL * 1.2, 0.2, (Math.random() - 0.5) * CELL * 1.2); grp.add(lump); }
    return grp;
  }
  function swarmModel() {
    var grp = new THREE.Group();
    for (var i = 0; i < 10; i++) { var b = new THREE.Mesh(new THREE.BoxGeometry(0.18, 0.1, 0.05), new THREE.MeshStandardMaterial({ color: 0x2a1d22, roughness: 1 })); b.position.set((Math.random() - 0.5) * 1.6, 1.1 + (Math.random() - 0.5) * 1.2, (Math.random() - 0.5) * 1.6); grp.add(b); }
    return grp;
  }
  function lureModel() {
    var grp = humanoid(0xc97b8a);
    var line = new THREE.Mesh(new THREE.CylinderGeometry(0.01, 0.01, 0.8, 6), new THREE.MeshStandardMaterial({ color: 0x222222 })); line.position.set(0.34, 2.1, 0); grp.add(line);
    var ball = new THREE.Mesh(new THREE.SphereGeometry(0.22, 12, 12), new THREE.MeshBasicMaterial({ color: 0xff3344 })); ball.position.set(0.34, 2.55, 0); grp.add(ball);
    return grp;
  }
  function captainClark() {
    var grp = humanoid(0x2e3a2a);
    var hat = new THREE.Mesh(new THREE.CylinderGeometry(0.3, 0.36, 0.18, 12), new THREE.MeshStandardMaterial({ color: 0x14110a, roughness: 1 })); hat.position.y = 1.96; grp.add(hat);
    return grp;
  }
  var TYPES = [
    { b: 'chase', speed: 3.0, sight: 36, reach: 1.0, model: function () { return humanoid(0xded8c8); } }, // faceling
    { b: 'chase', speed: 4.3, sight: 30, reach: 1.0, model: function () { return humanoid(0xcdbfa0); } }, // hound (fast)
    { b: 'still', speed: 6.0, sight: 70, reach: 1.0, model: function () { return humanoid(0xb6ad98); } }, // still-life
    { b: 'smiler', speed: 6.5, sight: 13, reach: 1.1, model: smilerModel }                                // smiler
  ];
  var ents = [];
  for (var ei = 0; ei < 8; ei++) {
    var rm = rooms[ri(2, rooms.length - 1)];
    var rc = roomCenter(rm), tp = TYPES[ri(0, TYPES.length - 1)];
    var mdl = tp.model(); mdl.position.set(rc.x, 0, rc.z); scene.add(mdl);
    ents.push({ x: rc.x, z: rc.z, mesh: mdl, t: tp, wx: rc.x, wz: rc.z, wt: 0 });
  }
  // Pass-2 specials
  function spawnSpecial(t) {
    var rm = rooms[ri(2, rooms.length - 1)], rc = roomCenter(rm);
    var mdl = t.model(); mdl.position.set(rc.x, 0, rc.z); scene.add(mdl);
    ents.push({ x: rc.x, z: rc.z, mesh: mdl, t: t, wx: rc.x, wz: rc.z, wt: 0, trig: false });
  }
  spawnSpecial({ b: 'zone', r: CELL * 0.9, model: function () { return hazardModel(0xff3322); } }); // Growth
  spawnSpecial({ b: 'zone', r: CELL * 0.9, model: function () { return hazardModel(0x33ff66); } }); // Green Glow
  spawnSpecial({ b: 'zone', r: CELL * 0.9, model: function () { return hazardModel(0x771111); } }); // Clumps
  spawnSpecial({ b: 'swarm', r: 2.3, speed: 2.2, model: swarmModel });                              // Deathmoths
  spawnSpecial({ b: 'lure', r: 1.1, speed: 7.5, trig: 9, model: lureModel });                       // Partygoer
  spawnSpecial({ b: 'chase', speed: 3.4, sight: 34, reach: 1.0, model: captainClark });             // Captain Clark

  function losTo(ex, ez, tx, tz, range) {
    var dx = tx - ex, dz = tz - ez, dist = Math.sqrt(dx * dx + dz * dz);
    if (dist > range) return false;
    var steps = Math.ceil(dist / (CELL * 0.5));
    for (var s = 1; s < steps; s++) {
      if (isWall(Math.floor((ex + dx * s / steps) / CELL), Math.floor((ez + dz * s / steps) / CELL))) return false;
    }
    return true;
  }
  function moveEnt(e, speed, dt, tx, tz) {
    var dx = tx - e.x, dz = tz - e.z, d = Math.sqrt(dx * dx + dz * dz) || 1;
    var nx = e.x + dx / d * speed * dt, nz = e.z + dz / d * speed * dt;
    if (!blocked(nx, e.z)) e.x = nx;
    if (!blocked(e.x, nz)) e.z = nz;
    e.mesh.position.set(e.x, 0, e.z);
    e.mesh.rotation.y = Math.atan2(px - e.x, pz - e.z);
  }
  function caught() {
    px = checkpoint.x; pz = checkpoint.z; inv = 2.2;
    if (caughtMsg) { caughtMsg.style.opacity = '1'; setTimeout(function () { caughtMsg.style.opacity = '0'; }, 900); }
  }
  function updateEntities(dt) {
    if (inv > 0) inv -= dt;
    // checkpoint = previous room visited
    var bi = -1, bd = 1e9;
    for (var i = 0; i < rooms.length; i++) { var rc = roomCenter(rooms[i]); var d = (rc.x - px) * (rc.x - px) + (rc.z - pz) * (rc.z - pz); if (d < bd) { bd = d; bi = i; } }
    if (bi !== lastRoomIdx) { if (lastRoomIdx >= 0) checkpoint = roomCenter(rooms[lastRoomIdx]); lastRoomIdx = bi; }

    var fx = -Math.sin(yaw), fz = -Math.cos(yaw);
    nearDist = 1e9;
    for (var k = 0; k < ents.length; k++) {
      var e = ents[k], dx = px - e.x, dz = pz - e.z, dist = Math.sqrt(dx * dx + dz * dz);
      if (dist < nearDist) nearDist = dist;
      if (hidden) { // tucked away — entities can't find you; they drift/idle
        if (e.t.b === 'chase' || e.t.b === 'swarm') { e.wt -= dt; if (e.wt <= 0) { e.wx = e.x + ri(-5, 5) * CELL; e.wz = e.z + ri(-5, 5) * CELL; e.wt = 3; } moveEnt(e, (e.t.speed || 2) * 0.4, dt, e.wx, e.wz); }
        continue;
      }
      var catchR = e.t.r != null ? e.t.r : (e.t.reach + 0.35);
      if (inv <= 0 && dist < catchR) { caught(); return; }
      if (e.t.b === 'chase') {
        if (losTo(e.x, e.z, px, pz, e.t.sight)) moveEnt(e, e.t.speed, dt, px, pz);
        else { e.wt -= dt; if (e.wt <= 0) { e.wx = e.x + ri(-6, 6) * CELL; e.wz = e.z + ri(-6, 6) * CELL; e.wt = 2.5; } moveEnt(e, e.t.speed * 0.4, dt, e.wx, e.wz); }
      } else if (e.t.b === 'still') {
        var dd = dist || 1, dot = (fx * (e.x - px) + fz * (e.z - pz)) / dd;   // >0 = in front
        var observed = dot > 0.55 && losTo(px, pz, e.x, e.z, e.t.sight);
        if (!observed) moveEnt(e, e.t.speed, dt, px, pz);
        else e.mesh.rotation.y = Math.atan2(dx, dz);
      } else if (e.t.b === 'smiler') { // lurk, charge when near — flees the flashlight
        var beam = flashOn && (fx * (e.x - px) + fz * (e.z - pz)) / (dist || 1) > 0.8 && dist < 20 && losTo(px, pz, e.x, e.z, 20);
        if (beam) moveEnt(e, e.t.speed * 0.9, dt, e.x + (e.x - px), e.z + (e.z - pz));
        else if (dist < e.t.sight) moveEnt(e, e.t.speed, dt, px, pz);
        else e.mesh.rotation.y = Math.atan2(dx, dz);
      } else if (e.t.b === 'swarm') { // drifting moths
        e.mesh.rotation.y += dt * 0.8;
        e.wt -= dt; if (e.wt <= 0) { e.wx = e.x + ri(-5, 5) * CELL; e.wz = e.z + ri(-5, 5) * CELL; e.wt = 3; }
        moveEnt(e, e.t.speed, dt, e.wx, e.wz);
      } else if (e.t.b === 'lure') { // benign until you get close, then charges
        if (!e.trig && dist < e.t.trig) e.trig = true;
        if (e.trig) moveEnt(e, e.t.speed, dt, px, pz);
        else e.mesh.rotation.y = Math.atan2(dx, dz);
      }
      // zone: static hazard — no movement, catch handled above
    }
  }

  // ── Objective: cyber-puzzle terminals gate the EXIT ──
  // Batch 1 puzzles. type: choice | multi | code. clue (optional) feeds the final code.
  var FINALCODE = 'B443ZERO13';
  var PUZZLES = [
    { id: 'P01', title: 'Password Cracker', type: 'code', q: 'Create a strong password (8+ chars: upper, lower, digit, symbol):', validate: function (v) { return v.length >= 8 && /[A-Z]/.test(v) && /[a-z]/.test(v) && /[0-9]/.test(v) && /[^A-Za-z0-9]/.test(v); }, clue: '1ST: B' },
    { id: 'P02', title: 'Phishing Email', type: 'choice', q: 'Which message is the phishing attempt?', o: ['IT: "reset your password at THIS LINK now"', 'HR: quarterly newsletter', 'Coworker calendar invite'], a: 0 },
    { id: 'P03', title: 'Firewall Switches', type: 'multi', q: 'OPEN only the safe ports, then submit.', items: [{ l: '23 — Telnet', pick: false }, { l: '443 — HTTPS', pick: true }, { l: '21 — FTP', pick: false }, { l: '80 — HTTP', pick: true }, { l: '3389 — RDP', pick: false }], clue: 'NEXT: 443' },
    { id: 'P04', title: 'Caesar Cipher', type: 'code', q: 'Decode (shift back by 1): "AFSP"', answer: 'ZERO', clue: 'THEN: ZERO' },
    { id: 'P05', title: 'MFA Factors', type: 'multi', q: 'Select the 3 VALID authentication factors, then submit.', items: [{ l: 'Password (know)', pick: true }, { l: 'Fingerprint (are)', pick: true }, { l: 'Phone code (have)', pick: true }, { l: 'Your birthday', pick: false }, { l: 'Lucky number', pick: false }] },
    { id: 'P06', title: 'Log Analysis', type: 'choice', q: 'Which login is suspicious?', o: ['09:14 — 10.0.0.5', '13:02 — 10.0.0.5', '03:47 — 10.0.6.13'], a: 2, clue: 'END: 13' },
    { id: 'P07', title: 'Hash Match', type: 'choice', q: 'Known-good hash is a1b2. Which file is TAMPERED?', o: ['report.docx → a1b2', 'budget.xlsx → a1b2', 'setup.exe → 9f3c'], a: 2 },
    { id: 'P08', title: 'Cable Maze', type: 'choice', q: 'Which path reaches the ROUTER without crossing the breach?', o: ['Server → SW1 → BREACH → Router', 'Server → SW2 → SW3 → Router', 'Server → BREACH → Router'], a: 1 },
    { id: 'P09', title: 'Social Engineering', type: 'choice', q: '"IT Support" calls demanding your password. You:', o: ['Give it — they are IT', 'Refuse and report it', 'Email it to be safe'], a: 1 },
    { id: 'P10', title: 'Malware Quarantine', type: 'multi', q: 'QUARANTINE only the unsafe files, then submit.', items: [{ l: 'invoice.pdf', pick: false }, { l: 'update.exe (unknown)', pick: true }, { l: 'photo.jpg', pick: false }, { l: 'free_vbucks.scr', pick: true }, { l: 'report.docx', pick: false }] },
    { id: 'P11', title: 'Encryption Keypad', type: 'code', q: 'Enter the recovered key:', answer: 'VAULT', clue: 'KEY VAULT' }
  ];
  var REG_CNT = Math.min(rooms.length - 1, PUZZLES.length);
  var NEED = REG_CNT, codes = 0, paused = false, won = false, activeTerm = null, clues = [];

  // One terminal per gated room (each opens the door to the next room when solved)
  var terms = [];
  for (var ti = 0; ti < REG_CNT; ti++) {
    var rc = roomCenter(rooms[ti]);
    var ped = new THREE.Mesh(new THREE.BoxGeometry(0.5, 1.1, 0.5), new THREE.MeshStandardMaterial({ color: 0x10202a, roughness: 0.7 }));
    ped.position.set(rc.x, 0.55, rc.z); scene.add(ped);
    var scr = new THREE.Mesh(new THREE.BoxGeometry(0.46, 0.34, 0.06), new THREE.MeshBasicMaterial({ color: 0x00e5ff }));
    scr.position.set(rc.x, 1.25, rc.z); scene.add(scr);
    terms.push({ x: rc.x, z: rc.z, scr: scr, p: PUZZLES[ti], solved: false, lock: 0, door: doors[ti] });
  }
  // P12 — final keypad in the last room (the EXIT)
  var p12ped = new THREE.Mesh(new THREE.BoxGeometry(0.5, 1.1, 0.5), new THREE.MeshStandardMaterial({ color: 0x2a1010, roughness: 0.7 })); p12ped.position.set(ex, 0.55, ez); scene.add(p12ped);
  var p12scr = new THREE.Mesh(new THREE.BoxGeometry(0.46, 0.34, 0.06), new THREE.MeshBasicMaterial({ color: 0xff7a2a })); p12scr.position.set(ex, 1.25, ez); scene.add(p12scr);
  terms.push({ x: ex, z: ez, scr: p12scr, p: { id: 'P12', title: 'Final Keypad', type: 'final', answer: FINALCODE }, solved: false, lock: 0 });

  // Locked door meshes (red while the gating room is unsolved)
  doors.forEach(function (d) {
    var ax = (d.cells[0][0] + d.cells[1][0]) / 2, az = (d.cells[0][1] + d.cells[1][1]) / 2;
    var horiz = d.cells[0][1] === d.cells[1][1];
    var mesh = new THREE.Mesh(new THREE.BoxGeometry(horiz ? CELL * 2.1 : CELL * 1.05, WALL_H, horiz ? CELL * 1.05 : CELL * 2.1),
      new THREE.MeshStandardMaterial({ color: 0xaa2222, roughness: 0.6, emissive: 0x330000 }));
    mesh.position.set(ax * CELL + CELL / 2, WALL_H / 2, az * CELL + CELL / 2);
    d.mesh = mesh; scene.add(mesh);
  });

  // Hiding spots (lockers) — in a few rooms, offset from the terminal
  var hideSpots = [];
  for (var hs = 0; hs < 4; hs++) {
    var hr = roomCenter(rooms[ri(1, rooms.length - 2)]);
    var locker = new THREE.Group();
    locker.add(new THREE.Mesh(new THREE.BoxGeometry(0.9, 2.2, 0.6), new THREE.MeshStandardMaterial({ color: 0x29382c, roughness: 0.8 })).translateY(1.1));
    locker.add(new THREE.Mesh(new THREE.BoxGeometry(0.05, 0.2, 0.06), new THREE.MeshStandardMaterial({ color: 0x9aa, roughness: 0.5 })).translateX(0.3).translateY(1.2).translateZ(0.3));
    var hx = hr.x + CELL * 2, hz = hr.z + CELL * 2;
    locker.position.set(hx, 0, hz); scene.add(locker);
    hideSpots.push({ x: hx, z: hz });
  }

  // EXIT barrier (red until unlocked)
  var barrier = new THREE.Mesh(new THREE.PlaneGeometry(CELL * 1.6, WALL_H),
    new THREE.MeshBasicMaterial({ color: 0xff2a2a, transparent: true, opacity: 0.35, side: THREE.DoubleSide }));
  barrier.position.set(ex, WALL_H / 2, ez); scene.add(barrier);

  // Overlay DOM
  var prompt2 = document.createElement('div');
  prompt2.style.cssText = 'position:fixed;top:54%;left:50%;transform:translateX(-50%);z-index:12;color:#00e5ff;font-family:monospace;font-size:18px;letter-spacing:2px;text-shadow:0 0 8px #000;display:none;';
  var obj = document.createElement('div');
  obj.style.cssText = 'position:fixed;top:14px;left:50%;transform:translateX(-50%);z-index:12;color:#e8dca0;font-family:monospace;font-size:15px;letter-spacing:1px;background:rgba(0,0,0,0.45);padding:6px 16px;border-radius:6px;pointer-events:none;';
  var clueTray = document.createElement('div');
  clueTray.style.cssText = 'position:fixed;top:14px;right:14px;z-index:12;color:#9effc8;font-family:monospace;font-size:13px;letter-spacing:1px;background:rgba(0,0,0,0.45);padding:8px 14px;border-radius:6px;pointer-events:none;display:none;';
  var puz = document.createElement('div');
  puz.style.cssText = 'position:fixed;inset:0;z-index:30;display:none;align-items:center;justify-content:center;background:rgba(2,6,10,0.92);font-family:monospace;';
  function updateObj() { obj.textContent = won ? 'ESCAPED' : ('TERMINALS ' + codes + '/' + NEED + (codes >= NEED ? ' · ENTER FINAL CODE AT EXIT' : '')); }
  function renderClues() {
    if (!clues.length) { clueTray.style.display = 'none'; return; }
    clueTray.style.display = 'block';
    clueTray.innerHTML = '<div style="color:#00e5ff;margin-bottom:4px;">CLUES</div>' + clues.map(function (c) { return '· ' + c; }).join('<br>');
  }

  var BTN = 'display:block;width:100%;margin:7px 0;padding:11px 15px;font:inherit;font-size:16px;color:#cfe;background:#0c1722;border:1px solid #1c3a4a;border-radius:8px;cursor:pointer;text-align:left;';
  function openDoor(d) { for (var i = 0; i < d.cells.length; i++) g[d.cells[i][1]][d.cells[i][0]] = 0; if (d.mesh) d.mesh.visible = false; }
  function solveTerm(term) {
    var p = term.p; term.solved = true; term.scr.material.color.setHex(0x16321f); codes++; updateObj();
    if (p.clue) { clues.push(p.clue); renderClues(); }
    if (term.door) openDoor(term.door);
    closePuzzle();
  }
  function failTerm(term) {
    var msg = puz.querySelector('#puz-msg'); if (msg) msg.textContent = 'ACCESS DENIED';
    if (term) term.lock = 2.0;
    setTimeout(closePuzzle, 750);
  }
  function winGame() {
    won = true; barrier.material.color.setHex(0x22ff66); barrier.material.opacity = 0.12;
    updateObj(); prompt2.style.display = 'none'; closePuzzle(); winMsg.style.display = 'flex'; document.exitPointerLock();
  }
  function openPuzzle(term) {
    paused = true; document.exitPointerLock();
    var p = term.p, h = '<div style="max-width:600px;width:90%;color:#cfe;text-align:center;">';
    h += '<div style="color:#00e5ff;font-size:13px;letter-spacing:3px;margin-bottom:8px;">' + p.id + ' · ' + p.title.toUpperCase() + '</div>';

    if (p.type === 'collect') {
      h += '<div style="font-size:19px;margin:8px 0 18px;">Collect 3 MFA tokens scattered in the rooms.</div>';
      h += '<div style="font-size:22px;color:#33ddff;">' + tokensGot + ' / 3 found</div>';
      h += '<button class="pzclose" style="' + BTN + 'text-align:center;margin-top:18px;">CLOSE</button></div>';
      puz.innerHTML = h; puz.style.display = 'flex';
      puz.querySelector('.pzclose').onclick = closePuzzle; return;
    }
    if (p.type === 'final') {
      if (codes < NEED) {
        h += '<div style="font-size:19px;margin:8px 0 18px;color:#ff7a5a;">INSUFFICIENT CLEARANCE</div>';
        h += '<div style="font-size:15px;color:#9bb;">Solve all terminals first (' + codes + '/' + NEED + ').</div>';
        h += '<button class="pzclose" style="' + BTN + 'text-align:center;margin-top:18px;">CLOSE</button></div>';
        puz.innerHTML = h; puz.style.display = 'flex'; puz.querySelector('.pzclose').onclick = closePuzzle; return;
      }
      h += '<div style="font-size:19px;margin:8px 0 18px;">Enter the final code (assemble it from your CLUES):</div>';
      h += '<input id="puz-input" autocomplete="off" style="width:70%;padding:12px;font:inherit;font-size:22px;text-align:center;letter-spacing:4px;text-transform:uppercase;color:#cfe;background:#0c1722;border:1px solid #1c3a4a;border-radius:8px;">';
      h += '<br><button class="pzsub" style="' + BTN + 'width:70%;margin:14px auto 0;text-align:center;color:#9effc8;border-color:#2a6a4a;">UNLOCK</button>';
      h += '<div id="puz-msg" style="height:20px;margin-top:12px;color:#ff5a5a;font-size:14px;"></div></div>';
      puz.innerHTML = h; puz.style.display = 'flex';
      var fin = puz.querySelector('#puz-input'); if (fin) fin.focus();
      puz.querySelector('.pzsub').onclick = function () {
        if ((fin.value || '').trim().toUpperCase().replace(/\s/g, '') === p.answer.toUpperCase()) winGame();
        else { var msg = puz.querySelector('#puz-msg'); if (msg) msg.textContent = 'CODE REJECTED'; }
      };
      return;
    }

    h += '<div style="font-size:19px;margin:8px 0 18px;">' + p.q + '</div>';
    if (p.type === 'choice') {
      for (var i = 0; i < p.o.length; i++) h += '<button class="pz" data-i="' + i + '" style="' + BTN + '">' + p.o[i] + '</button>';
    } else if (p.type === 'multi') {
      for (var j = 0; j < p.items.length; j++) h += '<button class="pzt" data-i="' + j + '" style="' + BTN + '">' + p.items[j].l + '</button>';
      h += '<button class="pzsub" style="' + BTN + 'text-align:center;color:#9effc8;border-color:#2a6a4a;margin-top:14px;">SUBMIT</button>';
    } else { // code
      h += '<input id="puz-input" autocomplete="off" style="width:60%;padding:12px;font:inherit;font-size:20px;text-align:center;letter-spacing:3px;color:#cfe;background:#0c1722;border:1px solid #1c3a4a;border-radius:8px;">';
      h += '<br><button class="pzsub" style="' + BTN + 'width:60%;margin:14px auto 0;text-align:center;color:#9effc8;border-color:#2a6a4a;">SUBMIT</button>';
    }
    h += '<div id="puz-msg" style="height:20px;margin-top:12px;color:#ff5a5a;font-size:14px;"></div></div>';
    puz.innerHTML = h; puz.style.display = 'flex';
    if (p.type === 'choice') {
      puz.querySelectorAll('.pz').forEach(function (b) { b.onclick = function () { if (parseInt(b.getAttribute('data-i'), 10) === p.a) solveTerm(term); else failTerm(term); }; });
    } else if (p.type === 'multi') {
      var sel = {};
      puz.querySelectorAll('.pzt').forEach(function (b) {
        b.onclick = function () { var k = b.getAttribute('data-i'); sel[k] = !sel[k]; b.style.background = sel[k] ? '#16321f' : '#0c1722'; b.style.borderColor = sel[k] ? '#2a6a4a' : '#1c3a4a'; };
      });
      puz.querySelector('.pzsub').onclick = function () {
        var ok = true; for (var k = 0; k < p.items.length; k++) if (!!sel[k] !== !!p.items[k].pick) ok = false;
        if (ok) solveTerm(term); else failTerm(term);
      };
    } else {
      var inp = puz.querySelector('#puz-input'); if (inp) inp.focus();
      puz.querySelector('.pzsub').onclick = function () {
        var v = (inp.value || '').trim();
        var ok = p.validate ? p.validate(v) : (v.toUpperCase() === p.answer.toUpperCase());
        if (ok) solveTerm(term); else failTerm(term);
      };
    }
  }
  function closePuzzle() { puz.style.display = 'none'; paused = false; }

  function checkInteract(fx, fz) {
    activeTerm = null; activeHide = null;
    var best = null, kind = null, bd = 2.4;
    for (var i = 0; i < terms.length; i++) {
      var t = terms[i]; if (t.solved || t.lock > 0) continue;
      var dx = t.x - px, dz = t.z - pz, d = Math.sqrt(dx * dx + dz * dz);
      if (d < bd && (fx * dx + fz * dz) / (d || 1) > 0.2) { bd = d; best = t; kind = 'term'; }
    }
    for (var h = 0; h < hideSpots.length; h++) {
      var s = hideSpots[h], hx = s.x - px, hz = s.z - pz, hd = Math.sqrt(hx * hx + hz * hz);
      if (hd < bd && (fx * hx + fz * hz) / (hd || 1) > 0.2) { bd = hd; best = s; kind = 'hide'; }
    }
    if (kind === 'term') { activeTerm = best; prompt2.textContent = '[E] ' + best.p.id + ' ' + best.p.title.toUpperCase(); prompt2.style.display = 'block'; }
    else if (kind === 'hide') { activeHide = best; prompt2.textContent = '[E] HIDE'; prompt2.style.display = 'block'; }
    else prompt2.style.display = 'none';
  }

  // ── Audio: fluorescent hum ──
  var ac = null, humOn = true;
  function startHum() {
    if (ac) return;
    try {
      ac = new (window.AudioContext || window.webkitAudioContext)();
      var g1 = ac.createGain(); g1.gain.value = 0.035; g1.connect(ac.destination);
      var o1 = ac.createOscillator(); o1.type = 'sawtooth'; o1.frequency.value = 120;
      var o2 = ac.createOscillator(); o2.type = 'sine'; o2.frequency.value = 60;
      var lp = ac.createBiquadFilter(); lp.type = 'lowpass'; lp.frequency.value = 900;
      o1.connect(lp); o2.connect(lp); lp.connect(g1); o1.start(); o2.start();
      window.__brGain = g1;
      var tg = ac.createGain(); tg.gain.value = 0; tg.connect(ac.destination);
      var to = ac.createOscillator(); to.type = 'sine'; to.frequency.value = 56; to.connect(tg); to.start();
      window.__brTension = tg;
    } catch (e) { }
  }
  function footstep() {
    if (!ac) return;
    var dur = 0.08, buf = ac.createBuffer(1, Math.floor(ac.sampleRate * dur), ac.sampleRate), d = buf.getChannelData(0);
    for (var i = 0; i < d.length; i++) d[i] = (Math.random() * 2 - 1) * Math.pow(1 - i / d.length, 3);
    var src = ac.createBufferSource(); src.buffer = buf;
    var bp = ac.createBiquadFilter(); bp.type = 'bandpass'; bp.frequency.value = 280;
    var g = ac.createGain(); g.gain.value = 0.22;
    src.connect(bp); bp.connect(g); g.connect(ac.destination); src.start();
  }

  // ── Input ──
  var keys = {};
  document.addEventListener('keydown', function (e) {
    if (paused || won) return;   // overlay open / escaped — let typing through, ignore game keys
    keys[e.code] = true;
    if (e.code === 'KeyM' && window.__brGain) { humOn = !humOn; window.__brGain.gain.value = humOn ? 0.035 : 0; }
    if (e.code === 'KeyF') flashOn = (!flashOn && battery > 3);
    if (e.code === 'KeyE') {
      if (hidden) hidden = false;
      else if (activeHide) { hidden = true; px = activeHide.x; pz = activeHide.z; }
      else if (activeTerm) openPuzzle(activeTerm);
    }
  });
  document.addEventListener('keyup', function (e) { keys[e.code] = false; });
  document.addEventListener('keydown', function (e) { if (e.code === 'Escape') window.location.href = '/page/start'; });
  renderer.domElement.addEventListener('click', function () { renderer.domElement.requestPointerLock(); startHum(); });
  document.addEventListener('mousemove', function (e) {
    if (document.pointerLockElement !== renderer.domElement) return;
    yaw -= e.movementX * 0.0024; pitch -= e.movementY * 0.0024;
    pitch = Math.max(-1.45, Math.min(1.45, pitch));
  });
  window.addEventListener('resize', function () {
    camera.aspect = window.innerWidth / window.innerHeight; camera.updateProjectionMatrix();
    renderer.setSize(window.innerWidth, window.innerHeight);
  });

  var hint = document.createElement('div');
  hint.style.cssText = 'position:fixed;bottom:16px;left:50%;transform:translateX(-50%);z-index:10;color:#e8dca0;font-family:monospace;font-size:13px;background:rgba(0,0,0,0.45);padding:6px 14px;border-radius:6px;pointer-events:none;letter-spacing:1px;';
  hint.textContent = 'CLICK LOOK · WASD MOVE · SHIFT SPRINT · F LIGHT · E USE · M MUTE';
  root.appendChild(hint);

  var caughtMsg = document.createElement('div');
  caughtMsg.style.cssText = 'position:fixed;inset:0;z-index:20;display:flex;align-items:center;justify-content:center;color:#ff3b3b;font-family:monospace;font-size:64px;font-weight:800;letter-spacing:4px;background:rgba(60,0,0,0.4);opacity:0;transition:opacity .25s;pointer-events:none;text-shadow:0 0 20px #000;';
  caughtMsg.textContent = 'CAUGHT';
  root.appendChild(caughtMsg);
  root.appendChild(prompt2); root.appendChild(obj); root.appendChild(puz); root.appendChild(clueTray);
  updateObj();

  var winMsg = document.createElement('div');
  winMsg.style.cssText = 'position:fixed;inset:0;z-index:40;display:none;align-items:center;justify-content:center;flex-direction:column;gap:20px;background:rgba(0,20,6,0.85);color:#22ff66;font-family:monospace;text-shadow:0 0 18px #000;';
  winMsg.innerHTML = '<div style="font-size:60px;font-weight:800;letter-spacing:6px;">YOU ESCAPED</div><div style="font-size:16px;color:#9fe;cursor:pointer;" onclick="location.reload()">[ click to re-enter ]</div>';
  root.appendChild(winMsg);

  var hud2 = document.createElement('div');
  hud2.style.cssText = 'position:fixed;bottom:14px;left:14px;z-index:12;font-family:monospace;font-size:11px;color:#cfe;letter-spacing:1px;pointer-events:none;';
  hud2.innerHTML = '<div>BATTERY</div><div style="width:120px;height:8px;background:#222;border:1px solid #444;margin:2px 0 6px;"><div id="br-bat" style="height:100%;width:100%;background:#ffd23f;"></div></div><div>STAMINA</div><div style="width:120px;height:8px;background:#222;border:1px solid #444;margin-top:2px;"><div id="br-sta" style="height:100%;width:100%;background:#33dd66;"></div></div>';
  root.appendChild(hud2);
  function updateHud() { var ba = document.getElementById('br-bat'), st = document.getElementById('br-sta'); if (ba) ba.style.width = battery + '%'; if (st) st.style.width = stamina + '%'; }

  var hideOv = document.createElement('div');
  hideOv.style.cssText = 'position:fixed;inset:0;z-index:14;display:none;align-items:flex-end;justify-content:center;padding-bottom:80px;background:radial-gradient(circle at center, rgba(0,0,0,0.3) 30%, rgba(0,0,0,0.92) 100%);color:#9effc8;font-family:monospace;font-size:18px;letter-spacing:2px;pointer-events:none;';
  hideOv.textContent = 'HIDDEN · [E] EXIT';
  root.appendChild(hideOv);

  var mm = document.createElement('canvas'); mm.width = mm.height = 150;
  mm.style.cssText = 'position:fixed;bottom:14px;right:14px;z-index:12;border:1px solid rgba(0,229,255,0.3);background:rgba(0,0,0,0.5);pointer-events:none;';
  root.appendChild(mm);
  var mmx = mm.getContext('2d'), MMS = 150 / Math.max(GW, GH);
  function dot(cx, cz, col, sz) { mmx.fillStyle = col; mmx.fillRect((cx / CELL) * MMS - sz, (cz / CELL) * MMS - sz, sz * 2, sz * 2); }
  function drawMinimap() {
    mmx.fillStyle = '#0a0f14'; mmx.fillRect(0, 0, 150, 150);
    mmx.fillStyle = '#1c2a33';
    for (var r = 0; r < GH; r++) for (var c = 0; c < GW; c++) if (g[r][c] === 0) mmx.fillRect(c * MMS, r * MMS, MMS + 0.5, MMS + 0.5);
    for (var i = 0; i < terms.length; i++) { var t = terms[i]; var col = t.p.type === 'final' ? (codes >= NEED ? '#22ff66' : '#ff7a2a') : (t.solved ? '#16321f' : '#00e5ff'); dot(t.x, t.z, col, 2); }
    dot(ex, ez, codes >= NEED ? '#22ff66' : '#ff7a2a', 2.5);
    dot(px, pz, '#ffd23f', 2.5);
  }

  // ── Loop ──
  var last = 0;
  function frame(t) {
    var dt = last ? Math.min(0.05, (t - last) / 1000) : 0.016; last = t;
    var fwx = -Math.sin(yaw), fwz = -Math.cos(yaw);
    var moving = false, sprintNow = false;
    if (!paused && !won && hidden) {
      updateEntities(dt);
    } else if (!paused && !won) {
      var fb = (keys.KeyW ? 1 : 0) - (keys.KeyS ? 1 : 0);
      var lr = (keys.KeyD ? 1 : 0) - (keys.KeyA ? 1 : 0);
      var sprinting = !!(keys.ShiftLeft || keys.ShiftRight) && stamina > 1 && (fb || lr);
      sprintNow = sprinting; moving = !!(fb || lr);
      stamina = sprinting ? Math.max(0, stamina - 26 * dt) : Math.min(100, stamina + 16 * dt);
      if (fb || lr) {
        var spd = SPEED * (sprinting ? 1.7 : 1);
        var rx = Math.cos(yaw), rz = -Math.sin(yaw);
        var mx = fwx * fb + rx * lr, mz = fwz * fb + rz * lr, ml = Math.sqrt(mx * mx + mz * mz) || 1;
        var dx = (mx / ml) * spd * dt, dz = (mz / ml) * spd * dt;
        if (!blocked(px + dx, pz)) px += dx;
        if (!blocked(px, pz + dz)) pz += dz;
      }
      updateEntities(dt);
      for (var ti = 0; ti < terms.length; ti++) if (terms[ti].lock > 0) terms[ti].lock -= dt;
      checkInteract(fwx, fwz);
    }
    // Flashlight (battery) + spotlight follow + HUD/minimap — always
    if (flashOn && battery > 0) { battery = Math.max(0, battery - 9 * dt); if (battery <= 0) flashOn = false; }
    else if (!flashOn) battery = Math.min(100, battery + 5 * dt);
    spot.intensity = flashOn ? 2.4 : 0;
    spot.position.set(px, EYE, pz);
    spotTarget.position.set(px + fwx * 5, EYE - 0.4, pz + fwz * 5);
    hideOv.style.display = hidden ? 'flex' : 'none';
    // Footsteps + head-bob
    if (moving) { bobPhase += dt * (sprintNow ? 11 : 7); stepT -= dt; if (stepT <= 0) { footstep(); stepT = sprintNow ? 0.32 : 0.48; } }
    else stepT = 0;
    var bob = moving ? Math.sin(bobPhase) * 0.06 : 0;
    // Fluorescent flicker
    flick -= dt; if (flick <= 0) { ambient.intensity = (Math.random() < 0.5) ? 0.32 : 0.72; flick = 0.05 + Math.random() * 0.28; }
    // Proximity tension
    if (window.__brTension) { var tv = hidden ? 0 : Math.max(0, Math.min(0.06, (10 - nearDist) / 10 * 0.06)); window.__brTension.gain.value = tv * (0.6 + 0.4 * Math.sin(t * 0.012)); }
    updateHud(); drawMinimap();
    camera.position.set(px, EYE + bob, pz); camera.rotation.y = yaw; camera.rotation.x = pitch;
    renderer.render(scene, camera);
    requestAnimationFrame(frame);
  }
  requestAnimationFrame(frame);
})();

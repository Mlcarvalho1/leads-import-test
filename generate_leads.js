const fs = require('fs');
const path = require('path');

const TOTAL_LEADS = 1000000;

const firstNames = [
  'João', 'Maria', 'Pedro', 'Ana', 'Carlos', 'Lucas',
  'Fernanda', 'Rafael', 'Juliana', 'Bruno', 'Camila'
];

const lastNames = [
  'Silva', 'Santos', 'Oliveira', 'Pereira', 'Costa',
  'Rodrigues', 'Alves', 'Lima', 'Gomes', 'Ribeiro'
];

const tagsPool = [
  'novo', 'vip', 'retorno', 'interessado',
  'premium', 'lead-frio', 'lead-quente'
];

function randomInt(min, max) {
  return Math.floor(Math.random() * (max - min + 1)) + min;
}

function randomPhone() {
  const ddd = randomInt(10, 99);
  const number = randomInt(10000000, 99999999);
  return `55${ddd}9${number}`;
}

function randomCPF() {
  return String(randomInt(0, 99999999999)).padStart(11, '0');
}

function randomName() {
  return `${firstNames[randomInt(0, firstNames.length - 1)]} ${lastNames[randomInt(0, lastNames.length - 1)]}`;
}

function randomEmail(name) {
  const base = name.toLowerCase().replace(/\s+/g, '.');
  return `${base}${randomInt(0, 999)}@email.com`;
}

function randomTags() {
  const count = randomInt(0, 5);
  if (count === 0) return '';

  const selected = new Set();
  while (selected.size < count) {
    selected.add(tagsPool[randomInt(0, tagsPool.length - 1)]);
  }

  return Array.from(selected).join(', ');
}

function generateCSV() {
  const start = Date.now();
  const filePath = path.join(__dirname, 'leads_10000.csv');
  const stream = fs.createWriteStream(filePath, { encoding: 'utf8' });

  // Header
  stream.write('name,phone,cpf,email,tags\n');

  for (let i = 0; i < TOTAL_LEADS; i++) {
    const name = randomName();
    const row = [
      name,
      randomPhone(),
      randomCPF(),
      randomEmail(name),
      randomTags()
    ].join(',');

    stream.write(row + '\n');
  }

  stream.end(() => {
    const elapsed = (Date.now() - start) / 1000;
    console.log(`✅ Arquivo leads_10000.csv gerado com sucesso! Tempo de execução: ${elapsed.toFixed(2)} segundos`);
  });
}

generateCSV();
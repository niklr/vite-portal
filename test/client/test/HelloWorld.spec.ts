import { describe } from "mocha";
import { expect } from "chai";
import * as vite from "@vite/vuilder";
import config from "./vite.config.json";
import { startRelayer } from "../src/vite"

let relayer: any;
let provider: any;
let deployer: any;

describe('test HelloWorld', () => {
  before(async function () {
    const relayerUrl = "http://127.0.0.1:56331";
    relayer = await startRelayer(relayerUrl);
    // const providerUrl = "http://127.0.0.1:23456";
    const providerUrl = relayerUrl + "/api/v1/client/relay";
    provider = vite.newProvider(providerUrl);
    deployer = vite.newAccount(config.networks.local.mnemonic, 0, provider);
    // console.log('deployer', deployer.address);
  });

  after(async function () {
    await relayer.stop();
  });

  it('test height', async function () {
    const height = await provider.request("ledger_getSnapshotChainHeight");
    expect(Number(height)).to.be.greaterThan(0);
  });

  xit('test contract', async () => {
    // compile
    const compiledContracts = await vite.compile('HelloWorld.solpp');
    expect(compiledContracts).to.have.property('HelloWorld');

    // deploy
    let helloWorld = compiledContracts.HelloWorld;
    helloWorld.setDeployer(deployer).setProvider(provider);
    await helloWorld.deploy({});
    expect(helloWorld.address).to.be.a('string');
    console.log(helloWorld.address);

    // check default value of data
    let result = await helloWorld.query('data', []);
    console.log('return', result);
    expect(result).to.be.an('array').with.lengthOf(1);
    expect(result![0]).to.be.equal('123');

    // call HelloWorld.set(456);
    await helloWorld.call('set', ['456'], {});

    // check value of data
    result = await helloWorld.query('data', []);
    console.log('return', result);
    expect(result).to.be.an('array').with.lengthOf(1);
    expect(result![0]).to.be.equal('456');
  });
});
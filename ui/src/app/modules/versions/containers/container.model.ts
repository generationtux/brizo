export class Container {
  public name: string;
  public image: string;
  public args: string[];
  public argString: string;
  public alwaysPull: boolean = true;
  public ports: ContainerPort[];
  public volumeMounts: VolumeMount[];

  constructor(props: any = {}) {
    Object.assign(this, props);
    this.args = this.args || [];
    this.ports = this.ports || [];
    this.volumeMounts = this.volumeMounts || [];
  }

  // convert args string into individual array of args
  public parseArgString() {
    if (this.args.length !== 0) {
      return;
    }

    let args = this.argString.split(",");
    let splitNoSpaces = args.map((value) => { return value.trim(); });

    let splitNoQuotes = splitNoSpaces.map((value) => {
      if (value.substring(0, 1) == '"') {
        value = value.substring(1);
      }
      if (value.substring(value.length-1) == '"') {
        value = value.substring(0, value.length-1);
      }

      return value;
    });

    this.args = splitNoQuotes;
  }
}

export class ContainerPort {
  public protocol: string;
  public port: number;

  constructor(props: any = {}) {
    Object.assign(this, props);
    this.protocol = this.protocol || "TCP";
  }
}

export class VolumeMount {
  public name: string;
  public path: string;

  constructor(props: any = {}) {
    Object.assign(this, props);
  }
}

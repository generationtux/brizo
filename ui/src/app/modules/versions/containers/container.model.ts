export class Container {
  public name: string;
  public image: string;
  public args: string[];
  public alwaysPull: boolean = true;
  public ports: ContainerPort[];
  public volumeMounts: VolumeMount[];

  constructor(props: any = {}) {
    Object.assign(this, props);
    this.args = this.args || [];
    this.ports = this.ports || [];
    this.volumeMounts = this.volumeMounts || [];
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

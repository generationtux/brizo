export class Version {
  public uuid: string;
  public name: string;
  public replicas: number;
  public slug: string;
  public image: string;
  public args: string[];
  public ports: ContainerPort[];
  public pullPolicy: string;
  public environment_uuid: string;
  public volumes: Volume[];
  public volumeMounts: VolumeMount[];

  constructor(props: any = {}) {
    Object.assign(this, props);
    this.pullPolicy = this.pullPolicy || "Always";
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

export class Volume {
    public name: string;
    public type: string;

    constructor(props: any = {}) {
      Object.assign(this, props);
      this.type = this.type || "temp";
    }
}

export class VolumeMount {
  public name: string;
  public path: string;

  constructor(props: any = {}) {
    Object.assign(this, props);
  }
}

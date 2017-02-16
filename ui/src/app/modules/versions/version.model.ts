export class Version {
  public uuid: string;
  public name: string;
  public replicas: number;
  public slug: string;
  public image: string;
  public args: string[];
  public ports: VersionPort[];
  public pullPolicy: string;
  public environment_uuid: string;

  constructor(props: any = {}) {
    Object.assign(this, props);
  }
}

export class VersionPort {
  public protocol: string;
  public port: number;

  constructor(props: any = {}) {
    Object.assign(this, props);
  }
}

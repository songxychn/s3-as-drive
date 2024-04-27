export namespace types {
	
	export class DownloadConfig {
	    dir: string;
	
	    static createFrom(source: any = {}) {
	        return new DownloadConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.dir = source["dir"];
	    }
	}
	export class S3Config {
	    endpoint: string;
	    accessKey: string;
	    secretKey: string;
	    bucket: string;
	
	    static createFrom(source: any = {}) {
	        return new S3Config(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.endpoint = source["endpoint"];
	        this.accessKey = source["accessKey"];
	        this.secretKey = source["secretKey"];
	        this.bucket = source["bucket"];
	    }
	}
	export class Config {
	    s3Config: S3Config;
	    downloadConfig: DownloadConfig;
	
	    static createFrom(source: any = {}) {
	        return new Config(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.s3Config = this.convertValues(source["s3Config"], S3Config);
	        this.downloadConfig = this.convertValues(source["downloadConfig"], DownloadConfig);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	export class Result {
	    code: number;
	    msg: string;
	    data: any;
	
	    static createFrom(source: any = {}) {
	        return new Result(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.code = source["code"];
	        this.msg = source["msg"];
	        this.data = source["data"];
	    }
	}

}


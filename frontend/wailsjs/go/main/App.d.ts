// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT
import {types} from '../models';

export function DeleteFile(arg1:number):Promise<void>;

export function DownloadFile(arg1:number):Promise<void>;

export function GetConfig():Promise<types.Config>;

export function GetFileList(arg1:string):Promise<Array<types.File>>;

export function GetShareUrl(arg1:number,arg2:number):Promise<void>;

export function Mkdir(arg1:string,arg2:string):Promise<void>;

export function UpdateConfig(arg1:types.Config):Promise<void>;

export function UploadDir(arg1:string):Promise<string>;

export function UploadFiles(arg1:string):Promise<number>;

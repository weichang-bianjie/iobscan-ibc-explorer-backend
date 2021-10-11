/* eslint-disable @typescript-eslint/camelcase */
import * as mongoose from 'mongoose';
import { IbcDenomType } from '../types/schemaTypes/ibc_denom.interface';
import { dateNow } from '../helper/date.helper';

export const IbcDenomSchema = new mongoose.Schema(
  {
    chain_id: String,
    denom: String,
    base_denom: String,
    denom_path: String,
    is_source_chain: Boolean,
    is_base_denom: {
      type: Boolean,
      default: false,
    },
    auth: {
      type: Boolean,
      default: false,
    },
    create_at: {
      type: Number,
      default: dateNow,
    },
    update_at: {
      type: Number,
      default: dateNow,
    },
    tx_time: {
      type: Number,
      default: dateNow,
    },
  },
  { versionKey: false },
);

IbcDenomSchema.index({ chain_id: 1, denom: 1 }, { unique: true });
IbcDenomSchema.index({ update_at: -1 }, { background: true });
IbcDenomSchema.index({ tx_time: -1 }, { background: true });

IbcDenomSchema.statics = {
  async findAllRecord(): Promise<IbcDenomType[]> {
    return this.find();
  },

  async findCount(): Promise<number> {
    return this.count();
  },

  async findBaseDenomCount(): Promise<number> {
    return this.count({
      is_source_chain: true
    });
  },

  async findDenomRecord(chain_id, denom): Promise<IbcDenomType> {
    return this.findOne({ chain_id, denom }, { _id: 0 });
  },

  async updateDenomRecord(denomRecord): Promise<void> {
    const { chain_id, denom } = denomRecord;
    const options = { upsert: true, new: false, setDefaultsOnInsert: true };
    return this.findOneAndUpdate({ chain_id, denom }, denomRecord, options);
  },

  async insertManyDenom(ibcDenom): Promise<void> {
    return this.insertMany(ibcDenom, { ordered: false });
  },
};
